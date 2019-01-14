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
	processorGroup goka.Group = "triptracker"
)

// stores the current trip status of a taxi
type taxiTrip struct {
	TaxiID    string `json:"taxi_id"`
	LicenseID string `json:"license_id"`

	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
}

type taxiTripsCodec int

func (ts *taxiTripsCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *taxiTripsCodec) Decode(data []byte) (interface{}, error) {
	var taxiTrips taxiTrip
	return &taxiTrips, json.Unmarshal(data, &taxiTrips)
}

func trackTrips(ctx goka.Context, msg interface{}) {

	var (
		taxiTrips *taxiTrip
	)
	t := ctx.Value()
	if t != nil {
		taxiTrips = t.(*taxiTrip)
	} else {
		taxiTrips = &taxiTrip{
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

	g := goka.DefineGroup(
		processorGroup,
		goka.Input(godays.TopicTripStarted, new(godays.TripStartedCodec), trackTrips),
		goka.Input(godays.TopicTripEnded, new(godays.TripEndedCodec), trackTrips),
		goka.Persist(new(taxiTripsCodec)),
	)

	proc, err := goka.NewProcessor(strings.Split(*brokers, ","), g)
	if err != nil {
		log.Fatalf("error creating trips processor: %v", err)
	}

	view, err := goka.NewView(strings.Split(*brokers, ","), goka.GroupTable(processorGroup), new(taxiTripsCodec))
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
			trips := val.(*taxiTrip)
			w.Write([]byte(fmt.Sprintf("%s -> busy: %t\n", id, !trips.Ended.IsZero())))
		}
	})

	go http.ListenAndServe(":8080", nil)
	go view.Run(context.Background())
	proc.Run(context.Background())
}
