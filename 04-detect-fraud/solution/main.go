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
	switch ev := msg.(type) {
	case *godays.TripStarted:
		ctx.Loopback(ev.LicenseID, &licenseAction{
			Ts:     ev.Ts,
			TaxiID: ev.TaxiID,
			Type:   pickup,
		})
	case *godays.TripEnded:
		ctx.Loopback(ev.LicenseID, &licenseAction{
			Ts:     ev.Ts,
			TaxiID: ev.TaxiID,
			Type:   dropoff,
		})
	}
}

func trackLicenses(ctx goka.Context, msg interface{}) {
	var lt *godays.LicenseTracker
	val := ctx.Value()
	if val == nil {
		lt = new(godays.LicenseTracker)
	} else {
		lt = val.(*godays.LicenseTracker)
	}
	if lt.Taxis == nil {
		lt.Taxis = make(map[string]bool)
	}
	loop := msg.(*licenseAction)
	switch loop.Type {
	case pickup:
		lt.Started = loop.Ts
		lt.Ended = time.Time{}
		lt.Taxis[loop.TaxiID] = true
		if len(lt.Taxis) > 1 {
			val := ctx.Join(goka.GroupTable(godays.LicenseConfigGroup))
			if val == nil || !val.(*godays.LicenseConfig).Fraud {
				ctx.Emit(godays.LicenseConfigTopic, ctx.Key(), &godays.LicenseConfig{Fraud: true})
			}
		}
	case dropoff:
		lt.Ended = loop.Ts
		delete(lt.Taxis, loop.TaxiID)
	default:
		log.Printf("invalid loop action type: %#v", loop)
	}

	ctx.SetValue(lt)
}
