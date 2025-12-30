package main

import (
	"Jetstream/config"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS
	fmt.Println("Connecting to NATS...")
	nc, err := nats.Connect("nats://nats-1:4222,nats://nats-2:4222,nats://nats-3:4222")
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Create a stream if it does not exist
	err = CreateStream(js)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(config.StreamName)

	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", config.StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     config.StreamName,
			Subjects: []string{config.StreamSubjects},
			Storage:  nats.FileStorage,
			Replicas: 3,
		})

		if err != nil {
			return err
		}
	}
	return nil
}
