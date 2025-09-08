package database

import (
	"miniMutualFund/models"

	"gorm.io/gorm"
)

type IOrderDB interface {
	Create(order *models.PlaceOrder) (*models.PlaceOrder, error)
	GetByID(orderID uint) (*models.PlaceOrder, error)
	Update(order *models.PlaceOrder) (*models.PlaceOrder, error)
	GetAll() ([]models.PlaceOrder, error)
}
type OrderDb struct {
	DB *gorm.DB
}

func NewOrderDB(db *gorm.DB) IOrderDB {
	return &OrderDb{db}
}

//	func (odb *OrderDb) Create(order *models.PlaceOrder) (*models.PlaceOrder, error) {
//		tx := odb.DB.Create(order)
//		if tx.Error != nil {
//			return nil, tx.Error
//		}
//		return order, nil
//	}
func (odb *OrderDb) Create(order *models.PlaceOrder) (*models.PlaceOrder, error) {
	err := odb.DB.Transaction(func(tx *gorm.DB) error {
		// 1️⃣ Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2️⃣ Update holdings
		var holding models.Holding
		result := tx.Where("order_id = ? AND scheme_code = ?", order.OrderId, order.SchemeCode).First(&holding)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// No existing holding, create new
				holding = models.Holding{
					OrderID:    int(order.OrderId),
					SchemeCode: order.SchemeCode,
					Units:      order.Units,
				}
				if err := tx.Create(&holding).Error; err != nil {
					return err
				}
			} else {
				return result.Error
			}
		} else {
			// Update existing holding
			switch order.TransType {
			case "B":
				holding.Units += order.Units
			case "S":
				holding.Units -= order.Units
			}
			if err := tx.Save(&holding).Error; err != nil {
				return err
			}
		}

		// 3️⃣ Ensure scheme exists
		var scheme models.Scheme
		res := tx.Where("scheme_code = ?", order.SchemeCode).First(&scheme)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				scheme = models.Scheme{
					SchemeCode: order.SchemeCode,
					SchemeName: order.SchemeCode, // can fetch real name from API
				}
				if err := tx.Create(&scheme).Error; err != nil {
					return err
				}
			} else {
				return res.Error
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}
func (odb *OrderDb) GetByID(id uint) (*models.PlaceOrder, error) {
	order := new(models.PlaceOrder)
	tx := odb.DB.First(order, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
func (odb *OrderDb) Update(order *models.PlaceOrder) (*models.PlaceOrder, error) {
	tx := odb.DB.Save(order) // Save will update if PK exists
	if tx.Error != nil {
		return nil, tx.Error
	}
	return order, nil
}
func (odb *OrderDb) GetAll() ([]models.PlaceOrder, error) {
	var orders []models.PlaceOrder
	tx := odb.DB.Find(&orders)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return orders, nil
}
