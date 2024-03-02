package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `gorm:"unique" json:"email"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	Role         string `json:"role"`
	Confirmed    bool   `json:"confirmed"`
	ConfirmToken string `json:"confirm_token"`
}
