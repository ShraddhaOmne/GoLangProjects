package handlers

import (
	"food-delivery/models"

	"gorm.io/gorm"
)

type IOrderEventDB interface {
	CreateOrder(order *models.OrdersEvent) (*models.OrdersEvent, error)
}
type OrderEventDb struct {
	DB *gorm.DB
}


