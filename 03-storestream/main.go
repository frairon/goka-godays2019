package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
	"github.com/spf13/pflag"
)

var (
	brokers = pflag.String("brokers", "localhost:9092", "brokers")
)

const (
	processorGroup goka.Group = "taxi-state"
)

type TaxiTrips struct {
	TaxiID    string `json:"taxi_id"`
	LicenseID string `json:"license_id"`

	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
}

type TaxiTripsCodec int

func (ts *TaxiTripsCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *TaxiTripsCodec) Decode(data []byte) (interface{}, error) {
	var taxiTrips TaxiTrips
	return &taxiTrips, json.Unmarshal(data, &taxiTrips)
}

type trips struct{}

func (tr *trips) consumeTrips(ctx goka.Context, msg interface{}) {

	var (
		taxiTrips *TaxiTrips
	)
	t := ctx.Value()
	if t != nil {
		taxiTrips = t.(*TaxiTrips)
	} else {
		taxiTrips = &TaxiTrips{
			TaxiID: ctx.Key(),
		}
	}

	switch ev := msg.(type) {
	case *godays.TripStarted:
		taxiTrips.LicenseID = ev.LicenseID
		taxiTrips.Started = ev.Ts
		taxiTrips.Ended = time.Time{}
	case *godays.TripEnded:
		taxiTrips.LicenseID = ev.LicenseID
		taxiTrips.Ended = ev.Ts
	default:
		ctx.Fail(fmt.Errorf("invalid message type: %T", msg))
	}

	ctx.SetValue(taxiTrips)
}

func main() {
	pflag.Parse()

	trp := new(trips)

	g := goka.DefineGroup(
		processorGroup,
		goka.Input(godays.TopicTripStarted, new(godays.TripStartedCodec), trp.consumeTrips),
		goka.Input(godays.TopicTripEnded, new(godays.TripEndedCodec), trp.consumeTrips),
		goka.Persist(new(TaxiTripsCodec)),
	)

	proc, err := goka.NewProcessor(strings.Split(*brokers, ","), g)
	if err != nil {
		log.Fatalf("error creating trips processor: %v", err)
	}

	view, err := goka.NewView(strings.Split(*brokers, ","), goka.GroupTable(processorGroup), new(TaxiTripsCodec))
	if err != nil {
		log.Fatalf("error creating trips view: %v", err)
	}

	http.HandleFunc("/taxi", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for _, id := range r.URL.Query()["id"] {
			val, _ := view.Get(id)
			if val == nil {
				w.Write([]byte(fmt.Sprintf("%s -> not found\n", id)))
				continue
			}
			trips := val.(*TaxiTrips)
			w.Write([]byte(fmt.Sprintf("%s -> busy: %t\n", id, !trips.Ended.IsZero())))
		}
	})

	go http.ListenAndServe(":8080", nil)
	go view.Run(context.Background())
	proc.Run(context.Background())
}
