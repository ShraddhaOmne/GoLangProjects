package models

type CommonModel struct {
	OrderId        uint `json:"id" gorm:"primaryKey;autoIncrement"`
	LastTradedTime int64
	Status         string
}
