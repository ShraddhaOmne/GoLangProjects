package database

import "users-service/models"

type IOrderDB interface {
	Create(order *models.Order) (*models.Order, error)
}
