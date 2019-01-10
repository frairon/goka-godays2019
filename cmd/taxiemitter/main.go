package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

const (
	readerChannelSize = 100
	timeFormat        = "2006-01-02 15:04:05"
)

var (
	input     = pflag.String("input", "", "input events file")
	timeLapse = pflag.Float64("time-lapse", 1.0, "increase or decrease time. >1.0 -> time runs faster")
)

func main() {
	pflag.Parse()

	f, err := os.Open(*input)
	if err != nil {
		log.Fatalf("Error opening file %s for reading: %v", *input, err)
	}

	defer f.Close()

	c := make(chan []string, 1000)

	reader := csv.NewReader(f)

	go func() {
		for {
			record, readErr := reader.Read()
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				log.Fatal(readErr)
			}
			c <- record
		}
		close(c)
	}()

	startTime := time.Now()

	emitEvent := func(eventTime time.Time, record []string) {
		log.Printf("Emitting %s", strings.Join(record, ", "))
	}

	firstEvent := <-c

	baseTime, err := time.Parse(timeFormat, firstEvent[0])
	if err != nil {
		log.Fatalf("Error parsing basetime %s: %v", firstEvent[0], err)
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

	wg.Wait()

}
