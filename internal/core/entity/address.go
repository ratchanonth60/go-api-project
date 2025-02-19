package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(64);not null"`
	Street      string `json:"street" gorm:"type:varchar(255);not null"`
	City        string `json:"city" gorm:"type:varchar(255);not null"`
	State       string `json:"state" gorm:"type:varchar(255);not null"`
	ZipCode     string `json:"zip_code" gorm:"type:varchar(255);not null"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(20);"`
	UserID      uint   `gorm:"not null;index"`
	User        User   `gorm:"foreignKey:UserID"`
}

func (a *Address) TableName() string {
	return "address"
}

func (a *Address) ToJson() ([]byte, error) {
	return json.Marshal(&a)
}
