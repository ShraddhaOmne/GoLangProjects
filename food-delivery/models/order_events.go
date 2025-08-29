package models

import "errors"

var (
	ErrInvalidStatus = errors.New("invalid status field")
)

type OrdersEvent struct {
	Id        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderId   uint   `json:"order_id" gorm:"not null;index"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Meta      string `json:"meta,omitempty"`
	Order     Orders `gorm:"foreignKey:OrderId`
}

func (oe *OrdersEvent) Validate() error {
	if oe.OrderId <= 0 {
		return ErrInvalidOrderId
	}
	if oe.Status == "" {
		return ErrInvalidStatus
	}
	return nil
}
