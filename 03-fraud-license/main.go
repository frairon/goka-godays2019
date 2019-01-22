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

			// <INSERT HERE>
			// configure input and persist
			// input stream is the godays.LicenseConfigTopic
			// persist as the same
		),
		utils.RandomStoragePath(),
	)
	if err != nil {
		log.Fatalf("Error creating processor: %v", err)
	}

	detectorProc, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(
			godays.LicenseDetectorGroup,

			// <INSERT HERE>
			// consume from trip-started
			// register lookup on the group table of processor 1
		),
		utils.RandomStoragePath(),
	)
	if err != nil {
		log.Fatalf("error creating view: %v", err)
	}

	// tip: start the processors using multierr
	me, ctx := multierr.NewErrGroup(context.Background())
	me.Go(func() error { return licenseConfProc.Run(ctx) })
	me.Go(func() error { return detectorProc.Run(ctx) })
	log.Fatal(me.Wait())
}

func configureLicense(ctx goka.Context, msg interface{}) {
	// <INSERT HERE>
	// store the config in the table
}

func detectBadLicenses(ctx goka.Context, msg interface{}) {
	// <INSERT HERE>
	// get the value from the lookup table and check if it's fraud.
}
