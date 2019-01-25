package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/multierr"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "brokers")
)

type action int

const (
	invalid action = iota
	pickup
	dropoff
)

type taxiCell struct {
	CurrentCell string `json:"current_cell"`
}

type cell struct {
	Lat          float64         `json:"lat"`
	Lon          float64         `json:"lon"`
	Updated      time.Time       `json:"updated"`
	TotalTip     float64         `json:"total_tip"`
	TaxisStarted map[string]bool `json:"taxis_started"`
	TaxisEnded   map[string]bool `json:"taxis_ended"`
}

// cellCodec encodes/decodes TripEnded
type cellCodec struct{}

// Encode encode message
func (ts *cellCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *cellCodec) Decode(data []byte) (interface{}, error) {
	var ca cell
	return &ca, json.Unmarshal(data, &ca)
}

// taxiCellCodec encodes/decodes taxiCell
type taxiCellCodec struct{}

// Encode encode message
func (ts *taxiCellCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *taxiCellCodec) Decode(data []byte) (interface{}, error) {
	var tc taxiCell
	return &tc, json.Unmarshal(data, &tc)
}

type cellAction struct {
	Lat          float64   `json:"lat"`
	Lon          float64   `json:"lon"`
	Ts           time.Time `json:"timestamp"`
	Tip          float64   `json:"tip"`
	PreviousCell string    `json:"prev_cell"`
	Action       action    `json:"action"`
	TaxiID       string    `json:"taxi_id"`
}

// cellActionCodec encodes/decodes TripEnded
type cellActionCodec struct{}

// Encode encode message
func (ts *cellActionCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes a message
func (ts *cellActionCodec) Decode(data []byte) (interface{}, error) {
	var ca cellAction
	return &ca, json.Unmarshal(data, &ca)
}

type uiCell struct {
	Lon             float64
	Lat             float64
	NumTaxisEnded   int
	NumTaxisStarted int
	TotalTip        float64
}

func main() {
	flag.Parse()

	proc, err := goka.NewProcessor(strings.Split(*brokers, ","), goka.DefineGroup(
		"taxi-cell",
		goka.Input(godays.TripStartedTopic, new(godays.TripStartedCodec), consumeTrips),
		goka.Input(godays.TripEndedTopic, new(godays.TripEndedCodec), consumeTrips),
		goka.Output("cell-action", new(cellActionCodec)),
		goka.Persist(new(taxiCellCodec)),
	))

	if err != nil {
		log.Fatalf("error creating trips processor: %v", err)
	}

	cellProc, err := goka.NewProcessor(strings.Split(*brokers, ","), goka.DefineGroup(
		"cells",
		goka.Input("cell-action", new(cellActionCodec), consumeCells),
		goka.Loop(new(cellActionCodec), removeTaxiFromPreviousCell),
		goka.Persist(new(cellCodec)),
	))
	if err != nil {
		log.Fatalf("error creating trips processor: %v", err)
	}

	view, err := goka.NewView(strings.Split(*brokers, ","),
		goka.GroupTable("cells"),
		new(cellCodec),
	)
	if err != nil {
		log.Fatalf("error creating trips view: %v", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(indexPage))
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		it, err := view.Iterator()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer it.Release()
		out := make(map[string]interface{})
		for it.Next() {
			v, _ := it.Value()
			c := v.(*cell)
			// ignore some out of order cells
			if c.Lon < -80 || c.Lon > -60 || c.Lat < 30 || c.Lat > 50 {
				continue
			}
			if len(c.TaxisStarted)+len(c.TaxisEnded) < 3 {
				continue
			}
			out[it.Key()] = &uiCell{
				Lat:             c.Lat,
				Lon:             c.Lon,
				NumTaxisEnded:   len(c.TaxisEnded),
				NumTaxisStarted: len(c.TaxisStarted),
				TotalTip:        c.TotalTip,
			}
		}
		log.Printf("found %d cells:", len(out))
		marshalled, err := json.Marshal(out)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Write(marshalled)
	})

	me, ctx := multierr.NewErrGroup(context.Background())
	me.Go(func() error { return proc.Run(ctx) })
	me.Go(func() error { return cellProc.Run(ctx) })
	me.Go(func() error { return view.Run(ctx) })
	me.Go(func() error {
		server := &http.Server{Addr: "localhost:8080"}
		go func() {
			server.ListenAndServe()
		}()
		defer server.Close()
		<-ctx.Done()
		return nil
	})
	log.Fatal(me.Wait())
}

func roundToCell(in float64) float64 {
	return math.Round(in*150.0) / 150.0
}

func toCell(lat, lon float64) string {
	return fmt.Sprintf("%.4f#%.4f", lat, lon)
}

func consumeTrips(ctx goka.Context, msg interface{}) {
	var tc *taxiCell
	if val := ctx.Value(); val != nil {
		tc = val.(*taxiCell)
	} else {
		tc = new(taxiCell)
	}
	switch ev := msg.(type) {
	case *godays.TripStarted:
		lat := roundToCell(ev.Latitude)
		lon := roundToCell(ev.Longitude)
		ctx.Emit("cell-action", toCell(lat, lon), &cellAction{
			Action:       pickup,
			PreviousCell: tc.CurrentCell,
			TaxiID:       ev.TaxiID,
			Ts:           ev.Ts,
			Lat:          lat,
			Lon:          lon,
		})

		// mark start cell of taxi
		tc.CurrentCell = toCell(lat, lon)
	case *godays.TripEnded:
		lat := roundToCell(ev.Latitude)
		lon := roundToCell(ev.Longitude)
		ctx.Emit("cell-action", toCell(lat, lon), &cellAction{
			Action:       dropoff,
			PreviousCell: tc.CurrentCell,
			Tip:          ev.Tip,
			TaxiID:       ev.TaxiID,
			Ts:           ev.Ts,
			Lat:          lat,
			Lon:          lon,
		})
		tc.CurrentCell = toCell(lat, lon)
	}

	ctx.SetValue(tc)
}
func consumeCells(ctx goka.Context, msg interface{}) {
	var c *cell

	if val := ctx.Value(); val != nil {
		c = val.(*cell)
	} else {
		c = new(cell)
	}
	if c.TaxisStarted == nil {
		c.TaxisStarted = make(map[string]bool)
	}
	if c.TaxisEnded == nil {
		c.TaxisEnded = make(map[string]bool)
	}
	ca := msg.(*cellAction)
	switch ca.Action {
	case pickup:
		c.TaxisStarted[ca.TaxiID] = true
		c.Lat = ca.Lat
		c.Lon = ca.Lon
		c.Updated = ca.Ts
	case dropoff:
		c.TaxisEnded[ca.TaxiID] = true
		c.Lat = ca.Lat
		c.Lon = ca.Lon
		c.Updated = ca.Ts
		c.TotalTip += ca.Tip
	}
	ctx.SetValue(c)
	ctx.Loopback(ca.PreviousCell, ca)
}

func removeTaxiFromPreviousCell(ctx goka.Context, msg interface{}) {
	ca := msg.(*cellAction)

	if val := ctx.Value(); val != nil {
		c := val.(*cell)
		if c.TaxisEnded != nil {
			delete(c.TaxisEnded, ca.TaxiID)
		}
		if c.TaxisStarted != nil {
			delete(c.TaxisStarted, ca.TaxiID)
		}
		ctx.SetValue(c)
	}
}

var indexPage = `
<html>

<head>
  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.4.0/dist/leaflet.css" integrity="sha512-puBpdR0798OZvTTbP4A8Ix/l+A4dHDD0DGqYW6RQ+9jxkRFclaxxQb/SJAWZfWAkuyeQUytO7+7N4QKrDh+drA=="
    crossorigin="" />
  <script src="https://unpkg.com/leaflet@1.4.0/dist/leaflet.js" integrity="sha512-QVftwZFqvtRNi0ZyCtsznlKSWOStnDORoefr1enyq5mVL4tmKB3S/EnC3rRJcxCPavG10IcrVGSmPh6Qw5lwrg=="
    crossorigin=""></script>
  <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>

  <style>
    #map {
      height: 1000px;
      width: 1000px;
    }
  </style>
</head>

<body>
  <div id="map"></div>

  <script type="text/javascript">
    var map = L.map('map');

    var osmUrl = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    var osmAttrib = 'Map data © <a href="https://openstreetmap.org">OpenStreetMap</a> contributors';
    var osm = new L.TileLayer(osmUrl, {
      minZoom: 8,
      maxZoom: 12,
      attribution: osmAttrib
    });

    map.setView(new L.LatLng(40.7400, -73.8800), 12);
    map.addLayer(osm);

    var dots = {};

    $(function() {
      window.setInterval(function() {
        jQuery.ajax({
          url: "http://localhost:8080/status",
          dataType: "json",
          success: function(data) {
            console.log(data);
            for (var cell in data) {
              var v = data[cell];

              if (dots.hasOwnProperty(cell)) {
                dots[cell].setRadius((v.NumTaxisStarted + v.NumTaxisEnded) * 8);
              } else {
                var circle = L.circle([v.Lat, v.Lon], {
                  color: 'red',
                  fillColor: '#f03',
                  fillOpacity: 0.5,
                  radius: (v.NumTaxisStarted + v.NumTaxisEnded) * 8
                }).addTo(map);
                dots[cell] = circle;
              }
            }

            // remove extra dots
            for (var cell in dots) {
              if (!data.hasOwnProperty(cell)) {
                dots[cell].remove(map);
              }
            }
          }
        });
      }, 1000);
    });
  </script>
</body>

</html>
`
