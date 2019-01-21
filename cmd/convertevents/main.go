package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

var (
	input  = flag.String("input", "", "raw event data file")
	output = flag.String("output", "", "path to splitted output file")
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type timedRecord struct {
	ts time.Time
	r  []string
}

// Original data layout:

// 00 medallion             an md5sum of the identifier of the taxi – vehicle bound
// 01 hack_license          an md5sum of the identifier for the taxi license
// 02 pickup_datetime       time when the passenger(s) were picked up
// 03 dropoff_datetime      time when the passenger(s) were dropped off
// 04 trip_time_in_secs     duration of the trip
// 05 trip_distance         trip distance in miles
// 06 pickup_longitude      longitude coordinate of the pickup location
// 07 pickup_latitude       latitude coordinate of the pickup location
// 08 dropoff_longitude     longitude coordinate of the drop-off location
// 09 dropoff_latitude      latitude coordinate of the drop-off location
// 10 payment_type          the payment method – credit card or cash
// 11 fare_amount           fare amount in dollars
// 12 surcharge             surcharge in dollars
// 13 mta_tax               tax in dollars
// 14 tip_amount            tip in dollars
// 15 tolls_amount          bridge and tunnel tolls in dollars
// 16 total_amount          total paid amount in dollars

func main() {

	flag.Parse()
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
				record[2], "pickup", taxiName, licenseName, record[7], record[6], "", "", "", "",
			},
		})
		records = append(records, &timedRecord{
			ts: endTs,
			r: []string{
				record[3], "dropoff", taxiName, licenseName, record[9], record[8], record[16], record[14], record[4], record[5],
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
