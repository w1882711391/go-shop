package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId      string `json:"user_id"`      //用户id
	PassWord    string `json:"pass_word"`    // 密码
	Gender      string `json:"gender"`       //性别
	PhoneNumber string `json:"phone_number"` // 电话
	Email       string `json:"email"`        //邮箱
	ImgPath     string `json:"img_path"`     // 头像
	IsMerchant  bool   `json:"is_merchant"`  //是否被拉黑
}
