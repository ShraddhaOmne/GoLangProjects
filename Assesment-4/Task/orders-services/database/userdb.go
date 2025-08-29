package database

import (
	"errors"
	"users-service/models"

	"gorm.io/gorm"
)

type IUserDB interface {
	Create(user *models.User) (*models.User, error)
	GetBy(id uint) (*models.User, error)
	GetByLimit(limit, offset int) ([]models.User, error)
	CreateOrder(order *models.Order) (*models.Order, error)
}
type UserDb struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) IUserDB {
	return &UserDb{db}
}

func (udb *UserDb) Create(user *models.User) (*models.User, error) {
	tx := udb.DB.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *UserDb) GetBy(id uint) (*models.User, error) {
	user := new(models.User)
	tx := udb.DB.Preload("Orders").First(user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *UserDb) GetByLimit(limit, offset int) ([]models.User, error) {
	var users []models.User
	tx := udb.DB.Limit(limit).Offset(offset).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func (udb *UserDb) CreateOrder(order *models.Order) (*models.Order, error) {
	_, err := udb.GetBy(order.UserId)
	if err != nil {
		return nil, errors.New("invalid user or user not found")
	}
	tx := udb.DB.Create(order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
