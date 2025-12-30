package main

import (
	"Jetstream/config"
	"Jetstream/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func publishReviews(js nats.JetStreamContext) {
	reviews, err := getReviews()
	if err != nil {
		log.Println(err)
		return
	}

	for _, oneReview := range reviews {

		// create random message intervals to slow down
		r := rand.Intn(1500)
		time.Sleep(time.Duration(r) * time.Millisecond)
		oneReview.SentTime = time.Now().UnixMilli() // record send time

		reviewString, err := json.Marshal(oneReview)
		if err != nil {
			log.Println(err)
			continue
		}

		// publish to REVIEWS.rateGiven subject
		// start := time.Now()
		_, err = js.Publish(config.SubjectNameReviewCreated, reviewString)
		if err != nil {
			log.Println(err)
		} else {
			// oneReview.SentTime = start
			log.Printf("Publisher  =>  Message: %s\n", oneReview.Text)
		}
	}
}

func getReviews() ([]models.Review, error) {
	rawReviews, err := os.ReadFile("./reviews.json")
	if err != nil {
		return nil, err
	}

	if len(rawReviews) == 0 {
		return nil, fmt.Errorf("reviews.json is empty")
	}

	var reviewsObj []models.Review
	er := json.Unmarshal(rawReviews, &reviewsObj)

	return reviewsObj, er
}
