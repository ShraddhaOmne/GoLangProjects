package models

type Holding struct {
	OrderID     int     `json:"order_id"`
	SchemeCode string  `json:"scheme_code"`
	Units      float64 `json:"units"`
}
