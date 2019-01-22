package main

import (
	"flag"
	"log"
	"strings"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers   = flag.String("brokers", "localhost:9092", "brokers")
	topic     = flag.String("topic", "", "topic to send messages to")
	separator = flag.String("separator", " ", "String to separate key from value")
	key       = flag.String("key", "", "Key to use to emit")
	value     = flag.String("value", "", "value to use to emit")
)

func main() {
	flag.Parse()

	if *topic == "" {
		log.Fatalf("topic must be set")
	}

	emitter, err := goka.NewEmitter(strings.Split(*brokers, ","), goka.Stream(*topic), new(codec.String))
	if err != nil {
		log.Fatalf("Error creating emitter: %v", err)
	}
	defer emitter.Finish()
	emitter.EmitSync(*key, *value)
}
