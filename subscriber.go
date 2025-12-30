package main

import (
	"Jetstream/config"
	"Jetstream/db"
	"Jetstream/models"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func subscribeReviews(js nats.JetStreamContext) {
	_, err := js.Subscribe(
		config.SubjectNameReviewCreated,
		func(m *nats.Msg) {

			var review models.Review
			err := json.Unmarshal(m.Data, &review)
			if err != nil {
				log.Println("JSON unmarshal failed:", err)
				return
			}

			// Calculate latency
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

			// Ack ONLY after successful DB write
			if err := m.Ack(); err != nil {
				log.Println("Ack failed:", err)
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
