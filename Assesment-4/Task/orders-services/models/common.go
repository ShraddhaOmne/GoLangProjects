package models

type CommonModel struct {
	Id           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Status       string `json:"status"`
	LastModified int64  `json:"last_modified" gorm:"index"`
}
