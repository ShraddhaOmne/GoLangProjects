package models

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidName   = errors.New("invalid name field")
	ErrInvalidEmail  = errors.New("invalid email field")
	ErrInvalidMobile = errors.New("invalid mobile field")
)

type User struct {
	CommonModel         // promoted field
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Mobile      string  `json:"mobile"`
	Orders      []Order `json:"orders,omitempty" gorm:"foreignKey:UserId"` //Orders can or cannot be used so omitempty
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidName
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Mobile == "" {
		return ErrInvalidMobile
	}
	return nil
}
func (u *User) ToBytes() []byte {
	bytes, _ := json.Marshal(u)
	return bytes
}
