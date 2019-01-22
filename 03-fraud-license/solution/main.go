package main

import (
	"context"
	"flag"
	"log"
	"strings"

	godays "github.com/frairon/goka-godays2019"
	"github.com/frairon/goka-godays2019/utils"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/multierr"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "brokers")
)

func main() {
	flag.Parse()

	// create bad licenses processor
	licenseConfProc, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(
			godays.LicenseConfigGroup,

			goka.Input(godays.LicenseConfigTopic, new(godays.LicenseConfigCodec), configureLicense),
			goka.Persist(new(godays.LicenseConfigCodec)),
		),
		utils.RandomStoragePath(),
	)
	if err != nil {
		log.Fatalf("Error creating processor: %v", err)
	}

	detectorProc, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(godays.LicenseDetectorGroup,

			goka.Input(godays.TripStartedTopic, new(godays.TripStartedCodec), detectBadLicenses),
			goka.Lookup(goka.GroupTable(godays.LicenseConfigGroup), new(godays.LicenseConfigCodec)),
		),
		utils.RandomStoragePath(),
	)
	if err != nil {
		log.Fatalf("error creating view: %v", err)
	}

	me, ctx := multierr.NewErrGroup(context.Background())
	me.Go(func() error { return licenseConfProc.Run(ctx) })
	me.Go(func() error { return detectorProc.Run(ctx) })
	log.Fatal(me.Wait())
}

func configureLicense(ctx goka.Context, msg interface{}) {
	// simply the message as state
	log.Printf("setting license config: %#v", msg.(*godays.LicenseConfig))
	ctx.SetValue(msg)
}

func detectBadLicenses(ctx goka.Context, msg interface{}) {
	// cast the incoming message, will not fail
	started := msg.(*godays.TripStarted)

	// lookup the license config and alert a fraud
	if val := ctx.Lookup(goka.GroupTable(godays.LicenseConfigGroup), started.LicenseID); val != nil {
		cfg := val.(*godays.LicenseConfig)
		if cfg.Fraud {
			log.Printf("Detected Taxi trip with blacklisted license: %#v", started)
		}
	}
}
