package main

import (
	"fmt"
	"time"
)

// OrderStatus represents possible statuses
type OrderStatus string

const (
	PLACED           OrderStatus = "PLACED"
	PREPARING        OrderStatus = "PREPARING"
	COOKING          OrderStatus = "COOKING"
	OUT_FOR_DELIVERY OrderStatus = "OUT_FOR_DELIVERY"
	DELIVERED        OrderStatus = "DELIVERED"
)

// OrderEvent records a status change
type OrderEvent struct {
	Status    OrderStatus
	Timestamp time.Time
}

// simulateOrder simulates automatic progression of an order
func simulateOrder(orderID int, eventChan chan<- OrderEvent) {
	statuses := []OrderStatus{
		PLACED,
		PREPARING,
		COOKING,
		OUT_FOR_DELIVERY,
		DELIVERED,
	}

	for _, status := range statuses {
		// record event
		event := OrderEvent{
			Status:    status,
			Timestamp: time.Now(),
		}
		eventChan <- event

		// simulate delay before next status
		time.Sleep(time.Second * 2) // configurable delay per step
	}

	close(eventChan) // close channel when finished
}

func main() {
	eventChan := make(chan OrderEvent)

	// start goroutine
	go simulateOrder(101, eventChan)

	// listen for events
	for event := range eventChan {
		fmt.Printf("Order Event: %-15s at %s\n", event.Status, event.Timestamp.Format(time.RFC3339))
	}
}
