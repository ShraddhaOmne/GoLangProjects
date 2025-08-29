package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type PriceFeed struct {
	Scrip string
	Price float64
}

func main() {
	scrips := []string{"AAPL", "GOOG", "TSLA"}
	fmt.Println("Generating ticks...")

	priceCh := make(chan PriceFeed)

	// Start goroutines for each ticker
	for _, ticker := range scrips {
		go generateLTP(ticker, priceCh)
	}

	timer := time.After(10 * time.Second)

out:
	for {
		select {
		case msg := <-priceCh:
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("%s :: %s -- %.2f\n", timestamp, msg.Scrip, msg.Price)

		case <-timer:
			fmt.Println("Time elapsed..")
			break out
		}
	}
}

func generateLTP(scrip string, ch chan<- PriceFeed) {
	for {
		price := rand.Float64() * 1000
		ch <- PriceFeed{Scrip: scrip, Price: price}
		time.Sleep(time.Second)
	}
}
