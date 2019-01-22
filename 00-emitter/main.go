package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	godays "github.com/frairon/goka-godays2019"
	"github.com/lovoo/goka"
)

const (
	readerChannelSize = 100
	timeFormat        = "2006-01-02 15:04:05"
)

var (
	brokers          = flag.String("brokers", "localhost:9092", "brokers")
	input            = flag.String("input", "testdata/taxidata_tiny.csv", "input events file")
	timeLapse        = flag.Float64("time-lapse", 60, "increase or decrease time. >1.0 -> time runs faster")
	licenseFraudRate = flag.Int("license-fraud-rate", 0.0, "Every nth license is a fraud license")
)

var (
	eventsSent int64
)

func main() {
	flag.Parse()

	f, err := os.Open(*input)
	if err != nil {
		log.Fatalf("Error opening file %s for reading: %v", *input, err)
	}

	defer f.Close()

	c := make(chan []string, 1000)

	// start the emitters
	startEmitter, err := goka.NewEmitter(strings.Split(*brokers, ","), godays.TopicTripStarted, new(godays.TripStartedCodec))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer startEmitter.Finish()
	endEmitter, err := goka.NewEmitter(strings.Split(*brokers, ","), godays.TopicTripEnded, new(godays.TripStartedCodec))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer endEmitter.Finish()

	// read from csv input file and send events to channel
	reader := csv.NewReader(f)
	var timeRead time.Time
	go func() {
		for {
			record, readErr := reader.Read()
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				log.Fatal(readErr)
			}
			eventTime, err := time.Parse(timeFormat, record[0])
			if err != nil {
				log.Fatalf("Error parsing event time %s: %v", record[0], err)
			}
			if timeRead.After(eventTime) {
				log.Printf("event misordering")
			}
			timeRead = eventTime

			c <- record
		}
		close(c)
	}()

	startTime := time.Now()

	firstEvent := <-c

	baseTime, err := time.Parse(timeFormat, firstEvent[0])
	if err != nil {
		log.Fatalf("Error parsing basetime %s: %v", firstEvent[0], err)
	}

	emitEvent := func(eventTime time.Time, record []string) {
		event := parseFromCsvRecord(baseTime, startTime, record)
		switch ev := event.(type) {
		case *godays.TripStarted:
			if ev.Latitude == 0 && ev.Longitude == 0 {
				return
			}
			startEmitter.Emit(ev.TaxiID, event)
		case *godays.TripEnded:
			if ev.Latitude == 0 && ev.Longitude == 0 {
				return
			}
			endEmitter.Emit(ev.TaxiID, event)
		default:
			log.Fatalf("unhandled event type: %v", event)
		}
		atomic.AddInt64(&eventsSent, 1)
	}

	// emit the first event now
	emitEvent(baseTime, firstEvent)

	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				record, ok := <-c
				if !ok {
					return
				}

				eventTime, err := time.Parse(timeFormat, record[0])
				if err != nil {
					log.Fatalf("Error parsing event time %s: %v", record[0], err)
				}

				realDiff := time.Since(startTime)
				eventDiff := time.Duration(float64(eventTime.Sub(baseTime)) / *timeLapse)

				// wait for the event to occur
				if eventDiff > realDiff {
					time.Sleep(eventDiff - realDiff)
				}

				emitEvent(eventTime, record)
			}
		}()
	}

	go printEventCounter()

	wg.Wait()
}

func parseFromCsvRecord(baseEventTime time.Time, baseTime time.Time, record []string) interface{} {
	eventTime, err := time.Parse(timeFormat, record[0])
	if err != nil {
		log.Fatalf("Error parsing event time %s: %v", record[0], err)
	}

	licenseSplit := strings.Split(record[3], "-")
	licenseNumber, err := strconv.ParseInt(licenseSplit[1], 10, 64)
	if err != nil {
		log.Fatalf("Error parsing license ID %s: %v", record[3], err)
	}

	if *licenseFraudRate > 0 {
		if licenseNumber > int64(*licenseFraudRate) && licenseNumber%int64(*licenseFraudRate) == 0 {
			licenseNumber = licenseNumber - int64(*licenseFraudRate)
			log.Printf("creating duplicate: %d", licenseNumber)
		}
	}

	licenseID := fmt.Sprintf("license-%d", licenseNumber)

	realEventTime := baseTime.Add(eventTime.Sub(baseEventTime))
	switch record[1] {
	case "pickup":
		return &godays.TripStarted{
			Ts:        realEventTime,
			TaxiID:    record[2],
			LicenseID: licenseID,
			Latitude:  mustParseFloat(record[4]),
			Longitude: mustParseFloat(record[5]),
		}
	case "dropoff":
		return &godays.TripEnded{
			Ts:        realEventTime,
			TaxiID:    record[2],
			LicenseID: licenseID,
			Latitude:  mustParseFloat(record[4]),
			Longitude: mustParseFloat(record[5]),
			Charge:    mustParseFloat(record[6]),
			Tip:       mustParseFloat(record[7]),
			Duration:  time.Duration(mustParseFloat(record[8]) * float64(time.Second)),
			Distance:  mustParseFloat(record[9]),
		}
	}
	log.Fatalf("Invalid record type: %#v", record)
	return nil
}

func mustParseFloat(strVal string) float64 {
	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		log.Fatalf("Error parsing strVal %s: %v", strVal, err)
	}
	return floatVal
}

func printEventCounter() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		log.Printf("sent %d events", atomic.LoadInt64(&eventsSent))
	}
}
