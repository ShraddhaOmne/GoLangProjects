// database/ordereventdb.go
package database

import (
	"food-delivery/models"
	"time"

	"gorm.io/gorm"
)

type IOrderEventDB interface {
	CreateEvent(orderID uint, status string) error
}

type OrderEventDb struct {
	DB *gorm.DB
}

func NewOrderEventDB(db *gorm.DB) IOrderEventDB {
	return &OrderEventDb{db}
}

func (edb *OrderEventDb) CreateEvent(orderID uint, status string) error {
	event := models.OrdersEvent{
		OrderId:   orderID,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}
	return edb.DB.Create(&event).Error
}
	