package model

import (
	"gorm.io/gorm"
	"time"
)

// Product 商品结构体
type Product struct {
	gorm.Model
	NickName      string `json:"nick_name"`      //商品名称
	OriginalPrice int64  `json:"original_price"` //商品原价
	SalePrice     int64  `json:"sale_price"`     //现价
	ImgPath       string `json:"img_path"`       //商品图片地址
	Content       string `json:"content"`        //商品描述信息
	State         bool   `json:"state"`          //0下架 1上架
	Sid           string `json:"sid"`            //卖家id
	SName         string `json:"SName"`          //卖家name
	Like          bool   `json:"like" gorm:"-"`  //是否收藏
	LikeNum       int64  `json:"like_num"`       //收藏数量
	IsOrder       bool   `json:"is_order"`       //是否购买
}

// CartItem 购物车中商品信息
type CartItem struct {
	gorm.Model
	ProductId uint   `json:"product_id"` //商品id CartItem.ProductId=Product.Id
	Sid       string `json:"sid"`        //收藏人id
	IsLike    bool   `json:"isLike"`     //是否收藏
}

// Order 订单表
type Order struct {
	OrderID      string    `json:"order_id"` // 订单ID
	Sid          string    `json:"sid"`      // 商家ID
	SName        string    `json:"s_name"`
	ProductName  string    `json:"product_name"`
	ProductID    string    `json:"product_id"`    // 商品ID
	OrderStatus  string    `json:"order_status"`  // 订单状态：待付款/已付款/已发货/已完成
	OrderTime    time.Time `json:"order_time"`    // 下单时间
	PaymentTime  time.Time `json:"payment_time"`  // 支付时间
	PayType      string    `json:"pay_type"`      //支付方式
	ShippingTime time.Time `json:"shipping_time"` // 发货时间
	OId          string    `json:"oid"`           //购买人id
	OName        string    `json:"o_name"`
	Address      string    `json:"address"` //配送地址
}
