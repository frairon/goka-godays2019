package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"github.com/spf13/pflag"
)

var (
	input  = pflag.String("input", "", "raw event data file")
	output = pflag.String("output", "", "path to splitted output file")
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type timedRecord struct {
	ts time.Time
	r  []string
}

func main() {

	pflag.Parse()
	file, err := os.Open(*input)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	r := csv.NewReader(file)

	output, err := os.Create(*output)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer output.Close()

	outputWriter := csv.NewWriter(output)
	defer outputWriter.Flush()

	taxiNames := make(map[string]string)
	licenseNames := make(map[string]string)

	var records []*timedRecord

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		startTs, err := time.Parse(timeFormat, record[2])
		if err != nil {
			log.Fatalf("invalid time format %s: %v", record[2], err)
		}
		endTs, err := time.Parse(timeFormat, record[3])
		if err != nil {
			log.Fatalf("invalid time format %s: %v", record[3], err)
		}

		taxiName := taxiNames[record[0]]
		if taxiName == "" {
			taxiName = fmt.Sprintf("taxi-%d", len(taxiNames))
		}
		taxiNames[record[0]] = taxiName

		licenseName := licenseNames[record[1]]
		if licenseName == "" {
			licenseName = fmt.Sprintf("license-%d", len(licenseNames))
		}
		licenseNames[record[1]] = licenseName

		records = append(records, &timedRecord{
			ts: startTs,
			r: []string{
				"pickup", taxiName, licenseName, record[2], record[7], record[6], "", "",
			},
		})
		records = append(records, &timedRecord{
			ts: endTs,
			r: []string{
				"dropoff", taxiName, licenseName, record[3], record[9], record[8], record[16], record[14],
			},
		})
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].ts.Before(records[j].ts)
	})

	for _, record := range records {
		outputWriter.Write(record.r)
	}
}
