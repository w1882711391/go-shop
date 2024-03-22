package model

import (
	"gorm.io/gorm"
	"time"
)

// Product 商品结构体
type Product struct {
	gorm.Model
	NickName      string    `json:"nick_name"`      //商品名称
	OriginalPrice int64     `json:"original_price"` //商品原价
	SalePrice     int64     `json:"sale_price"`     //现价
	ImgPath       string    `json:"img_path"`       //商品图片地址
	Content       string    `json:"content"`        //商品描述信息
	State         bool      `json:"state"`          //0下架 1上架
	PubTime       time.Time `json:"pub_time"`       //上架时间
	UserId        string    `json:"user_id"`        //卖家id
}

// CartItem 购物车中商品信息
type CartItem struct {
	gorm.Model
	ProductId uint   `json:"product_id"` //商品id CartItem.ProductId=Product.Id
	NickName  string `json:"nick_name"`  //商品名称
	SalePrice int64  `json:"sale_price"` //现价
	ImgPath   string `json:"img_path"`   //商品图片地址
	UserId    string `json:"user_id"`    //收藏id
}

type WantProduct struct {
	gorm.Model
	NickName    string    `json:"nick_name"`   //商品名称
	Content     string    `json:"content"`     //商品描述信息
	Price       int64     `json:"price"`       //求购价格
	UserId      string    `json:"user_id"`     //求购id
	Comment     string    `json:"comment"`     //评论信息
	CommentTime time.Time `json:"commentTime"` //评论时间
}

// Order 订单表
type Order struct {
	OrderID      int64     // 订单ID
	UserID       int64     // 用户ID
	ProductID    int64     // 商品ID
	OrderStatus  string    // 订单状态：待付款/已付款/已发货/已完成
	OrderTime    time.Time // 下单时间
	PaymentTime  time.Time // 支付时间
	PayType      string    //支付方式
	ShippingTime time.Time // 发货时间
}

// Comments 评论表
type Comments struct {
	CommentID      int       // 评论ID
	ProductID      int       // 商品ID
	UserID         int       // 用户ID
	CommentContent string    // 评论内容
	CommentTime    time.Time // 评论时间
}
