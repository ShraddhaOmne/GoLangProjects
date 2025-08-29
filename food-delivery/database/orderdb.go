package database

import (
	"fmt"
	"food-delivery/models"
	"time"

	"gorm.io/gorm"
)

type IOrderDB interface {
	Create(order *models.Orders) (*models.Orders, error)
	UpdateStatus(orderID uint, status string) error
	GetWithEvents(orderID uint) (*models.Orders, error)
}
type OrderDb struct {
	DB *gorm.DB
}

func NewOrderDB(db *gorm.DB) IOrderDB {
	return &OrderDb{db}
}

func (odb *OrderDb) Create(order *models.Orders) (*models.Orders, error) {
	tx := odb.DB.Create(order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
func (odb *OrderDb) UpdateStatus(orderID uint, status string) error {
	fmt.Printf("Order id :%d Status :%s\n", orderID, status)
	tx := odb.DB.Model(&models.Orders{}).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now().Unix(),
		})
	return tx.Error
}
func (odb *OrderDb) GetWithEvents(orderID uint) (*models.Orders, error) {
	order := new(models.Orders)
	tx := odb.DB.Preload("Events").First(order, orderID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
