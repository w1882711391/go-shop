package main

import (
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"go_shop/route"
)

func main() {
	if err := dao.MysqlInit(); err != nil {
		panic("数据库启动失败")
	}
	logrus.Info("数据库已启动")
	model.AutoMigrate()
	logrus.Info("数据库表创建成功")
	app := route.RouterInit()
	logrus.Info("服务器已启动，监听端口8080")

	if err := app.Listen(":8080"); err != nil {
		return
	}
}
