package main

import (
	"context"
	"flag"
	"log"
	"strings"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "brokers")
)

func trackBadLicenses(ctx goka.Context, msg interface{}) {
	log.Printf("bad license tracked: %s -> %s", ctx.Key(), msg)
	ctx.SetValue(msg)
}

func detectBadLicenses(ctx goka.Context, msg interface{}) {
	started := msg.(*godays.TripStarted)
	badLicense := ctx.Lookup(goka.GroupTable(godays.BadLicenseGroup), started.LicenseID)
	if badLicense != nil {
		blocked := badLicense.(string)
		if blocked == "blacklisted" {
			log.Printf("Detected Taxi trip with blacklisted license: %#v", started)
		}
	}
}

func main() {
	flag.Parse()

	badLicenseProc, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(
			godays.BadLicenseGroup,
			goka.Input(godays.ConfigureLicenseTopic, new(codec.String), trackBadLicenses),
			goka.Persist(new(codec.String)),
		))
	if err != nil {
		log.Fatalf("Error creating processor: %v", err)
	}

	detector, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(godays.LicenseDetectorGroup,
			goka.Input(godays.TopicTripStarted, new(godays.TripStartedCodec), detectBadLicenses),
			goka.Lookup(goka.GroupTable(godays.BadLicenseGroup), new(codec.String)),
		))
	if err != nil {
		log.Fatalf("error creating view: %v", err)
	}

	go badLicenseProc.Run(context.Background())
	detector.Run(context.Background())
}
