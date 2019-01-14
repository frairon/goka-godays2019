package main

import (
	"context"
	"log"
	"strings"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/spf13/pflag"
)

var (
	brokers = pflag.String("brokers", "localhost:9092", "brokers")
)

const (
	configureLicenseTopic goka.Stream = "configure-licenses"
	badLicenseGroup       goka.Group  = "bad-licenses"
	licenseDetector       goka.Group  = "license-detector"
)

func trackBadLicenses(ctx goka.Context, msg interface{}) {
	log.Printf("bad license tracked: %s -> %s", ctx.Key(), msg)
	ctx.SetValue(msg)
}

func detectBadLicenses(ctx goka.Context, msg interface{}) {
	started := msg.(*godays.TripStarted)
	badLicense := ctx.Lookup(goka.GroupTable(badLicenseGroup), started.LicenseID)
	if badLicense != nil {
		blocked := badLicense.(string)
		if blocked == "blacklisted" {
			log.Printf("Detected Taxi trip with blacklisted license: %#v", started)
		}
	}
}

func main() {
	pflag.Parse()

	badLicenseProc, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(
			badLicenseGroup,
			goka.Input(configureLicenseTopic, new(codec.String), trackBadLicenses),
			goka.Persist(new(codec.String)),
		))
	if err != nil {
		log.Fatalf("Error creating processor: %v", err)
	}

	detector, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(licenseDetector,
			goka.Input(godays.TopicTripStarted, new(godays.TripStartedCodec), detectBadLicenses),
			goka.Lookup(goka.GroupTable(badLicenseGroup), new(codec.String)),
		))
	if err != nil {
		log.Fatalf("error creating view: %v", err)
	}

	go badLicenseProc.Run(context.Background())
	detector.Run(context.Background())
}
