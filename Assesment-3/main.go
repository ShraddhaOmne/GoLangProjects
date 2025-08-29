package main

import (
	"fmt"
)

func main() {
	o1 := MarketOrder{"INFY", 50}
	o2 := LimitOrder{"AAPL", 100, 174.25, 173.50}
	o3 := LimitOrder{"GOOG", 100, 120.00, 125.00}

	processOrder(o1)
	processOrder(o2)
	processOrder(o3)
}

type Order interface {
	Execute() error
}

type MarketOrder struct {
	symbol string
	qty    int
}

func (mkt MarketOrder) Execute() error {
	fmt.Printf("Processing Market Order : Buying %d %s @ market price\n", mkt.qty, mkt.symbol)
	return nil
}

type LimitOrder struct {
	symbol       string
	qty          int
	limitPrice   float64
	currentPrice float64
}

func (limit LimitOrder) Execute() error {
	if limit.limitPrice < limit.currentPrice {
		return fmt.Errorf("limit price %0.2f is below market price %0.2f cannot place order for %s", limit.limitPrice, limit.currentPrice,limit.symbol)
	}
	fmt.Printf("Processing Limit Order: Buying %d %s @ %0.2f\n", limit.qty, limit.symbol, limit.limitPrice)
	return nil
}

func processOrder(o Order) {
	if err := o.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
