package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId      string `json:"user_id"`
	PassWord    string `json:"pass_word"`
	PhoneNumber string `json:"phone_number"`
	IsMerchant  bool   `json:"is_merchant"`
}
