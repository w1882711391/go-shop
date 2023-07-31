package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"go_shop/route"
	"os/signal"
	"syscall"
)

var listenHttpErrChan = make(chan error)

func init() {
	dao.MysqlInit()
	logrus.Info("mysql数据库已启动")
	dao.RedisInit()
	logrus.Info("redis数据库已启动")
	model.AutoMigrate()
	logrus.Info("数据库表创建成功")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		app := route.RouterInit()
		listenHttpErrChan <- app.Listen(":8080")
	}()

	select {
	case err := <-listenHttpErrChan:
		logrus.Errorf("http err: %+v\n", err)
	case <-ctx.Done():
		logrus.Info("Shutting down gracefully...")
	}
}
