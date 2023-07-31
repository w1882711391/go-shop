package model

import "gorm.io/gorm"

// Product 商品结构体
type Product struct {
	gorm.Model
	NickName    string `json:"nick_name"`   //商品名称
	Price       int16  `json:"price"`       //商品价格
	Description string `json:"description"` //商品描述信息
	Stock       int16  `json:"stock"`       //库存数量
	UserId      string `json:"user_id"`     //卖家id
}

// CartItem 购物车中商品信息
type CartItem struct {
	Product
	UserId string `json:"user_id"`
	Number int16  `json:"number"` //购物车中商品的数量
}
