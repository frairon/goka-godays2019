package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/frairon/goka-godays2019/utils"
	"github.com/lovoo/goka"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "brokers")
)

func main() {
	flag.Parse()

	// define our processor
	g := goka.DefineGroup(
		godays.TripTrackerGroup,
		// input topics
		goka.Input(godays.TopicTripStarted, new(godays.TripStartedCodec), consumeEvents),
		goka.Input(godays.TopicTripEnded, new(godays.TripEndedCodec), consumeEvents),
		// state codec
		goka.Persist(new(godays.TaxiStatusCodec)),
	)

	// create the processor
	proc, err := goka.NewProcessor(strings.Split(*brokers, ","), g,
		// use random storage path to avoid clashes (just for this workshop)
		utils.RandomStoragePath(),
	)
	if err != nil {
		log.Fatalf("error creating trips processor: %v", err)
	}

	view, err := goka.NewView(strings.Split(*brokers, ","),
		goka.GroupTable(godays.TripTrackerGroup),
		new(godays.TaxiStatusCodec),
		utils.RandomStorageViewPath(),
	)
	if err != nil {
		log.Fatalf("error creating trips view: %v", err)
	}

	http.HandleFunc("/taxi", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for _, id := range r.URL.Query()["id"] {

			// Try to get the taxi status from the View
			val, _ := view.Get(id)
			if val == nil {
				w.Write([]byte(fmt.Sprintf("%s -> not found\n", id)))
				continue
			}
			trips := val.(*godays.TaxiStatus)
			w.Write([]byte(fmt.Sprintf("%s > Trips: %d, busy duration %.2f minutes, pause duration %.2f minutes, busy: %t\n",
				id,
				trips.NumTrips,
				trips.BusyDuration.Minutes(),
				trips.PauseDuration.Minutes(),
				!trips.Ended.IsZero())))
		}
	})

	go http.ListenAndServe(":8080", nil)

	go view.Run(context.Background())

	proc.Run(context.Background())
}

func consumeEvents(ctx goka.Context, msg interface{}) {

	var (
		taxiStatus *godays.TaxiStatus
	)
	t := ctx.Value()
	if t != nil {
		taxiStatus = t.(*godays.TaxiStatus)
	} else {
		taxiStatus = &godays.TaxiStatus{
			TaxiID: ctx.Key(),
		}
	}

	switch ev := msg.(type) {
	case *godays.TripStarted:
		taxiStatus.LicenseID = ev.LicenseID
		taxiStatus.Started = ev.Ts

		taxiStatus.NumTrips++
		// if we came from an earlier trip, track the pause time
		if !taxiStatus.Ended.IsZero() {
			taxiStatus.PauseDuration = taxiStatus.PauseDuration + ev.Ts.Sub(taxiStatus.Started)
		}
		taxiStatus.Ended = time.Time{}
	case *godays.TripEnded:
		taxiStatus.LicenseID = ev.LicenseID
		taxiStatus.Ended = ev.Ts
		taxiStatus.BusyDuration += ev.Duration
	default:
		ctx.Fail(fmt.Errorf("invalid message type: %T", msg))
	}

	ctx.SetValue(taxiStatus)
}
