package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName           string `json:"user_name" gorm:"type:varchar(100);not null;uniqueIndex"`
	FirstName          string `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName           string `json:"last_name" gorm:"type:varchar(100);not null"`
	Email              string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Password           string `json:"-" gorm:"type:varchar(255);not null"`
	Identity           string `json:"identity" gorm:"type:varchar(20);not null;uniqueIndex"`
	IsActive           bool   `json:"is_active" gorm:"default:false"`
	ConfirmToken       string `gorm:"unique"`
	ResetPasswordToken string `gorm:"type:varchar(255)"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) ToJson() ([]byte, error) {
	return json.Marshal(&u)
}
