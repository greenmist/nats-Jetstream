package main

import (
	"Jetstream/config"
	"Jetstream/db"
	"Jetstream/models"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

func subscribeReviews(js nats.JetStreamContext) {
	_, err := js.Subscribe(config.SubjectNameReviewCreated, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			log.Println("Unable to Ack", err)
			return
		}

		var review models.Review
		err = json.Unmarshal(m.Data, &review)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Subscriber  =>  Subject: %s  -  ID: %s  -  Author: %s  -  Rating: %d\n", m.Subject, review.Id, review.Author, review.Rating)

		err = db.InsertReview(review.Id, review.Author, review.Store, review.Text, review.Rating)
		if err != nil {
			log.Println("Failed to insert into DB:", err)
		}

	})

	if err != nil {
		log.Println("Subscribe failed")
		return
	}
}
