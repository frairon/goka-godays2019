package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/frairon/goka-godays2019/utils"
	"github.com/lovoo/goka"
	""
)

var (
	brokers = String("brokers", "localhost:9092", "brokers")
)

func trackTrips(ctx goka.Context, msg interface{}) {

	var (
		taxiTrips *godays.TaxiTrip
	)
	t := ctx.Value()
	if t != nil {
		taxiTrips = t.(*godays.TaxiTrip)
	} else {
		taxiTrips = &godays.TaxiTrip{
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
	Parse()

	g := goka.DefineGroup(
		godays.TripTrackerGroup,
		goka.Input(godays.TopicTripStarted, new(godays.TripStartedCodec), trackTrips),
		goka.Input(godays.TopicTripEnded, new(godays.TripEndedCodec), trackTrips),
		goka.Persist(new(godays.TaxiTripsCodec)),
	)

	proc, err := goka.NewProcessor(strings.Split(*brokers, ","), g,
		utils.RandomStoragePath(),
	)
	if err != nil {
		log.Fatalf("error creating trips processor: %v", err)
	}

	view, err := goka.NewView(strings.Split(*brokers, ","),
		goka.GroupTable(godays.TripTrackerGroup),
		new(godays.TaxiTripsCodec),
		utils.RandomStorageViewPath(),
	)
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
			trips := val.(*godays.TaxiTrip)
			w.Write([]byte(fmt.Sprintf("%s -> busy: %t\n", id, !trips.Ended.IsZero())))
		}
	})

	go http.ListenAndServe(":8080", nil)
	go view.Run(context.Background())
	proc.Run(context.Background())
}
