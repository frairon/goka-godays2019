package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"strings"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "brokers")
)

const (
	licenseTracker goka.Group = "license-tracker"
)

type actionType int

const (
	invalid actionType = iota
	pickup
	dropoff
)

// internal loop message, should only be used by this processor
type licenseAction struct {
	Ts time.Time

	TaxiID string

	// type of the event
	Type actionType
}

type licenseActionCodec struct{}

func (ts *licenseActionCodec) Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (ts *licenseActionCodec) Decode(data []byte) (interface{}, error) {
	var la licenseAction
	return &la, json.Unmarshal(data, &la)
}

func main() {
	flag.Parse()
	detector, err := goka.NewProcessor(strings.Split(*brokers, ","),
		goka.DefineGroup(licenseTracker,
			goka.Input(godays.TripStartedTopic, new(godays.TripStartedCodec), processTrips),
			goka.Input(godays.TripEndedTopic, new(godays.TripEndedCodec), processTrips),
			goka.Loop(new(licenseActionCodec), trackLicenses),
			goka.Output(godays.LicenseConfigTopic, new(godays.LicenseConfigCodec)),
			goka.Persist(new(godays.LicenseTrackerCodec)),
			goka.Join(goka.GroupTable(godays.LicenseConfigGroup), new(godays.LicenseConfigCodec)),
		))
	if err != nil {
		log.Fatalf("error creating view: %v", err)
	}

	detector.Run(context.Background())
}

func processTrips(ctx goka.Context, msg interface{}) {
	// <INSERT HERE>
	// create a licenseAction message and do a loopback using
	// the license-ID
}

func trackLicenses(ctx goka.Context, msg interface{}) {
	// <INSERT HERE>
	// (1) load the table-value and create it if it's not set.
	// (2) depending on the licenseAction-Type, mark the currently active taxi
	// (3) If there are multiple taxis active, Emit a LicenseConfig configuring Fraud to
	//     to the LicenseConfig topic
	// (4) to avoid sending the fraud multiple times, join with the LicenseConfigGroup-table
	//     and check whether the license is already configured as fraud.
	// (5) don't forget to save the table-value
}
