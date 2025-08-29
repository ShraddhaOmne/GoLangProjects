package models

import "errors"

type Order struct {
	CommonModel
	UserId   uint    `json:"user_id"`
	Amount   float32 `json:"amount"`
	ItemName string  `json:"item_name" gorm:"column:item_name"`
}

func (o *Order) Validate() error {
	if o.UserId <= 0 {
		return errors.New("invalid user id")
	}
	if o.ItemName == "" {
		return errors.New("invalid Item Name")
	}
	return nil
}
