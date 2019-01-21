package main

// TODO: implement with sarama-cluster
import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	godays "github.com/frairon/goka-godays2019"
	"github.com/spf13/pflag"
)

var (
	brokers = pflag.String("brokers", "localhost:9092", "comma separated list of brokers")
)

func main() {
	pflag.Parse()
	client, err := sarama.NewClient(strings.Split(*brokers, ","), sarama.NewConfig())
	if err != nil {
		log.Fatalf("Error creating sarama client: %v", err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()
	partitionConsumer, err := consumer.ConsumePartition(string(godays.TopicTripStarted), 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()
	for {
		select {
		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				return
			}
			var tripStarted godays.TripStarted
			json.Unmarshal(msg.Value, &tripStarted)
			log.Printf("Partition 0: Taxi %s started.", tripStarted.TaxiID)
		case e, ok := <-partitionConsumer.Errors():
			if !ok {
				return
			}
			log.Printf("Received error: %v", e.Err)
		}
	}
}
