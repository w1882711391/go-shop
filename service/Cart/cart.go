package Cart

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"time"
)

// Cart 购物车
type Cart struct {
}

// AddItem 向购物车中添加商品
func (c *Cart) AddItem(ctm model.CartItem) error {
	//如果原来有这个商品的话
	if IsCart(ctm.NickName) {
		var cart model.CartItem
		err := dao.DB.Where("nick_name=?", ctm.NickName).First(&cart).Error
		if err != nil {
			return fmt.Errorf("数据库查询错误: %v", err)
		}
		totalNumber := cart.Number + ctm.Number
		if totalNumber >= 20 {
			return errors.New("超出最大可添加数量")
		}
		//开启一个事务 防止操作不一致
		tx := dao.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		cart.Number = totalNumber
		ctm.UpdatedAt = time.Now()

		err = dao.DB.Save(&cart).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("数据库保存错误: %v", err)
		}
		// 提交事务
		err = tx.Commit().Error
		if err != nil {
			return fmt.Errorf("事务提交错误: %v", err)
		}
	} else {
		ctm.UpdatedAt = time.Now()
		err := dao.DB.Create(&ctm).Error
		if err != nil {
			return fmt.Errorf("数据库保存错误: %v", err)
		}
	}

	logrus.Info("添加商品成功")
	return nil
}

func IsCart(nickname string) bool {
	var cart model.CartItem
	dao.DB.Where("nick_name=?", nickname).Find(&cart)
	if cart.NickName == nickname {
		return true
	}
	return false
}
