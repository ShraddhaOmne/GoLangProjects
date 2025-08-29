package database

import (
	"assesment-2/models"
	"errors"

	"gorm.io/gorm"
)

type ITradeDB interface {
	Create(user *models.Trade) (*models.Trade, error)
	GetBy(id uint) (*models.Trade, error)
	//GetPositions(order *models.Trade) (*models.Trade, error)
}
type TradeDb struct {
	DB *gorm.DB
}

func NewTradeDB(db *gorm.DB) ITradeDB {
	return &TradeDb{db}
}

func (udb *TradeDb) Create(user *models.Trade) (*models.Trade, error) {
	tx := udb.DB.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *TradeDb) GetBy(id uint) (*models.Trade, error) {
	user := new(models.Trade)
	tx := udb.DB.Preload("Trades").First(user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *TradeDb) CreateOrder(order *models.Trade) (*models.Trade, error) {
	_, err := udb.GetBy(order.OrderId)
	if err != nil {
		return nil, errors.New("invalid user or user not found")
	}
	tx := udb.DB.Create(order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
