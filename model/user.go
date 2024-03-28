package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId      string `json:"sid" form:"sid" gorm:"column:sid;"`                //用户id
	PassWord    string `json:"password" form:"password" gorm:"column:password;"` // 密码
	UserName    string `json:"user_name" form:"user_name" gorm:"column:user_name"`
	Gender      string `json:"gender"`       //性别
	PhoneNumber string `json:"phone_number"` // 电话
	Email       string `json:"email"`        //邮箱
	ImgPath     string `json:"img_path"`     // 头像
	IsMerchant  bool   `json:"is_merchant"`  //是否被拉黑
}
