package model

import (
	"github.com/sirupsen/logrus"
	"go_shop/dao"
)

func AutoMigrate() {

	if err := dao.DB.AutoMigrate(&Product{}); err != nil {
		logrus.Error("商品表自动迁移失败")
	}
	if err := dao.DB.AutoMigrate(&CartItem{}); err != nil {
		logrus.Error("购物车表自动迁移失败")
	}
	if err := dao.DB.AutoMigrate(&User{}); err != nil {
		logrus.Error("用户表自动迁移失败")
	}
	if err := dao.DB.AutoMigrate(&Order{}); err != nil {
		logrus.Error("订单表自动迁移失败")
	}
	return
}
