package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "brokers")
)

const (
	licenseTracker goka.Group = "license-tracker"
)

// internal loop message, that's why not exported
type licenseAction struct {
	Ts     time.Time
	Type   string
	TaxiID string
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
			goka.Input(godays.TripStartedTopic, new(godays.TripStartedCodec), remapLicenses),
			goka.Input(godays.TripEndedTopic, new(godays.TripEndedCodec), remapLicenses),
			goka.Loop(new(licenseActionCodec), trackLicenses),
			goka.Output("configure-licenses", new(codec.String)),
			goka.Persist(new(godays.LicenseTrackerCodec)),
		))
	if err != nil {
		log.Fatalf("error creating view: %v", err)
	}

	detector.Run(context.Background())
}

func remapLicenses(ctx goka.Context, msg interface{}) {
	switch ev := msg.(type) {
	case *godays.TripStarted:
		ctx.Loopback(ev.LicenseID, &licenseAction{
			Ts:     ev.Ts,
			TaxiID: ev.TaxiID,
			Type:   "started",
		})
	case *godays.TripEnded:
		ctx.Loopback(ev.LicenseID, &licenseAction{
			Ts:     ev.Ts,
			TaxiID: ev.TaxiID,
			Type:   "ended",
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
	switch ev := msg.(type) {
	case *licenseAction:
		if ev.Type == "started" {
			if !lt.Fraud && len(lt.Taxis) > 1 {
				ctx.Emit("configure-licenses", ctx.Key(), "blacklisted")
				lt.Fraud = true
			}
			lt.Started = ev.Ts
			lt.Ended = time.Time{}
			lt.Taxis[ev.TaxiID] = true
		} else {
			lt.Ended = ev.Ts
			delete(lt.Taxis, ev.TaxiID)
		}
	default:
		ctx.Fail(fmt.Errorf("Expected licenseAction in loop, got %T", msg))
	}

	ctx.SetValue(lt)
}
