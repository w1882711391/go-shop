package dao

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MysqlInit() error {
	dsn := "root:123456@tcp(localhost:3306)/xiaowang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接失败")
	}
	DB = db
	return nil
}
