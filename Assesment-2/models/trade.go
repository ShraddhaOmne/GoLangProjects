package models

import (
	"errors"
)

type Trade struct {
	CommonModel
	Symbol    string
	TransType string
	Quantity  int
	Price     float64
}

func (o *Trade) Validate() error {
	if o.OrderId <= 0 {
		return errors.New("invalid ID")
	}
	if o.Symbol == "" {
		return errors.New("invalid Item Name")
	}
	if o.TransType == ""{
		return errors.New("invalid Transaction Type")
	}
	if o.Quantity == 0 {
		return errors.New("invalid Quantity")
	}
	if o.Price == 0.0 {
		return errors.New("invalid Price")
	}
	return nil
}
