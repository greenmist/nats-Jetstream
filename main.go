package main

import (
	"Jetstream/db"
	"fmt"
	"log"
	"sync"
)

func main() {

	fmt.Println("Starting...")

	db.InitDB()

	js, err := JetStreamInit()
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		publishReviews(js)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		subscribeReviews(js)
	}()

	wg.Wait()

}
