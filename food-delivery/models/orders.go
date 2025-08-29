package models

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidOrderId      = errors.New("invalid order_id field")
	ErrInvalidCustomerName = errors.New("invalid customer_name field")
	ErrInvalidAddress      = errors.New("invalid address field")
	ErrInvalidItem         = errors.New("invalid item field")
	ErrInvalidSize         = errors.New("invalid size field")
)

type Orders struct {
	OrderId      uint          `json:"order_id" gorm:"primaryKey;autoIncrement"`
	CustomerName string        `json:"customer_name"`
	Address      string        `json:"address"`
	Item         string        `json:"item"`
	Size         string        `json:"size"`
	Status       string        `json:"status"`
	CreatedAt    int64         `json:"created_at"`
	UpdatedAt    int64         `json:"updated_at"`
	Events       []OrdersEvent `json:"events" gorm:"foreignKey:OrderId"`
}

func (o *Orders) Validate() error {
	if o.CustomerName == "" {
		return ErrInvalidCustomerName
	}
	if o.Address == "" {
		return ErrInvalidStatus
	}
	if o.Item == "" {
		return ErrInvalidItem
	}
	if o.Size == "" {
		return ErrInvalidSize
	}
	if o.Status == "" {
		return ErrInvalidStatus
	}
	return nil
}
func (u *Orders) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}
