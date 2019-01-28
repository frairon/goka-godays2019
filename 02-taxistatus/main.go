package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

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

		// <INSERT HERE>
		// input topics
		// persist
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
			_ = id
			// <INSERT HERE>
			// get the state from the view
			// print the information you want to see
		}
	})

	go http.ListenAndServe(":8080", nil)

	go view.Run(context.Background())

	// <INSERT HERE>
	// start the processor (blocking)
	_ = proc
}

func consumeEvents(ctx goka.Context, msg interface{}) {

	// <INSERT HERE>

	// (1) get the table value from context

	// (2) create it if nil or cast it if non-nil

	// (3) switch on the message types (remember we're getting TripStarted AND TripEnded in this callback)
	// on either message, decide what you want to track to the TaxiStatus state

	// (4) persist the state again
}
