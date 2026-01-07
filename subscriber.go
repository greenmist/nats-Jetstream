package main

import (
	"Jetstream/config"
	"Jetstream/db"
	"Jetstream/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func subscribeReviews(js nats.JetStreamContext) {
	start := time.Now()
	count := 0

	_, err := js.Subscribe(
		config.SubjectNameReviewCreated,
		func(m *nats.Msg) {

			var review models.Review
			err := json.Unmarshal(m.Data, &review)
			if err != nil {
				log.Println("JSON unmarshal failed:", err)
				return
			}

			//latency
			receivedTime := time.Now().UnixMilli()
			latency := float64(receivedTime - review.SentTime)
			log.Printf(
				"Subscriber => Subject: %s | ID: %s | Author: %s | Rating: %d | Latency: %v\n",
				m.Subject,
				review.Id,
				review.Author,
				review.Rating,
				latency,
			)

			err = db.InsertReview(
				review.Id,
				review.Author,
				review.Store,
				review.Text,
				review.Rating,
				latency,
			)
			if err != nil {
				log.Println("DB insert failed:", err)
				return
			}

			if err := m.Ack(); err != nil {
				log.Println("Ack failed:", err)
			}

			//throughput
			count++
			elapsed := time.Since(start).Seconds()
			if elapsed >= 10 {
				throughput := float64(count) / elapsed
				log.Println("Ack----------------------------:", throughput)
				fmt.Printf("Throughput: %.2f messages/sec\n", throughput)
				count = 0
				start = time.Now()
			}
		},
		nats.Durable("review-consumer"),
		nats.ManualAck(),
	)

	if err != nil {
		log.Println("Subscribe failed:", err)
		return
	}
}
