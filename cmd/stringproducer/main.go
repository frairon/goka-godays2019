package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var (
	brokers   = flag.String("brokers", "localhost:9092", "brokers")
	topic     = flag.String("topic", "", "topic to send messages to")
	separator = flag.String("separator", " ", "String to separate key from value")
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
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Got error: %v. Terminating.", err)
			break
		}
		text = strings.TrimSpace(text)
		splitted := strings.Split(text, *separator)
		if text == "" || len(splitted) < 1 {
			continue
		}
		var (
			key   string
			value string
		)
		key = splitted[0]
		if len(splitted) > 1 {
			value = splitted[1]
		}
		log.Printf("emitting %s: %s -> %s", key, value, *topic)
		emitter.Emit(key, value)
	}
}
