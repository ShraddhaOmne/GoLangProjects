package models

import (
	"encoding/json"
	"errors"
)

type PlaceOrder struct {
	OrderId    uint    `json:"order_id" gorm:"primaryKey;autoIncrement"`
	SchemeCode string  `json:"scheme_code"`
	TransType  string  `json:"trans_type"` // B or S
	Amount     float64 `json:"amount"`     // investment amount
	CreatedAt  int64   `json:"created_at"`
	Status     string  `json:"status"`
	Units      float64 `json:"units"`
	NAVUsed    float64 `json:"nav_used"`
}

func (o *PlaceOrder) Validate() error {
	if o.SchemeCode == "" {
		return errors.New("invalid Scheme Code")
	}
	if o.TransType != "B" && o.TransType != "S" {
		return errors.New("invalid Transaction Type")
	}
	if o.Amount < 0.0 {
		return errors.New("invalid Amount")
	}
	return nil
}
func (u *PlaceOrder) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}
