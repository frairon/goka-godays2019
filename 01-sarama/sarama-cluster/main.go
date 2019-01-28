package main

import (
	"encoding/json"
	"flag"
	"log"
	"strings"

	cluster "github.com/bsm/sarama-cluster"
	godays "github.com/frairon/goka-godays2019"
)

var (
	brokers = flag.String("brokers", "localhost:9092", "comma separated list of brokers")
)

func main() {
	flag.Parse()
	client, err := cluster.NewClient(strings.Split(*brokers, ","), cluster.NewConfig())
	if err != nil {
		log.Fatalf("Error creating sarama client: %v", err)
	}
	defer client.Close()

	consumer, err := cluster.NewConsumerFromClient(client, "sarama-consumer", []string{string(godays.TripStartedTopic)})
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				return
			}
			var tripStarted godays.TripStarted
			json.Unmarshal(msg.Value, &tripStarted)
			log.Printf("Partition %d: Taxi %s started.", msg.Partition, tripStarted.TaxiID)
		case e, ok := <-consumer.Errors():
			if !ok {
				return
			}
			log.Printf("Received error: %v", e.Error())
		}
	}
}
