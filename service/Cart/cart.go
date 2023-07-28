package Cart

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"time"
)

// AddItem 向购物车中添加商品
func AddItem(ctx *fiber.Ctx, ctm model.CartItem) error {

	userID, st := ctx.Locals("user_id").(string)
	if !st {
		return errors.New("类型断言失败")
	}
	// 类型断言成功，可以将 strUserID 作为字符串类型使用
	ctm.UserId = userID

	//如果原来有这个商品的话
	if IsCart(ctm.NickName) {
		var cart model.CartItem
		err := dao.DB.Table("cart_items").Where("user_id=? AND nick_name=?", userID, ctm.NickName).First(&cart).Error
		if err != nil {
			return fmt.Errorf("数据库查询错误: %v", err)
		}
		totalNumber := cart.Number + ctm.Number
		// 查询商品库存
		var p model.Product
		err = dao.DB.Table("products").Where("nick_name=?", ctm.NickName).First(&p).Error
		if err != nil {
			return fmt.Errorf("数据库查询错误: %v", err)
		}

		if totalNumber >= p.Stock {
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
			tx.Rollback()
			return fmt.Errorf("事务提交错误: %v", err)
		}
	} else {
		ctm.UpdatedAt = time.Now()
		err := dao.DB.Table("cart_items").Create(&ctm).Error
		if err != nil {
			return fmt.Errorf("数据库保存错误: %v", err)
		}
	}

	logrus.Info("添加商品至购物车成功")
	return nil
}

func IsCart(nickname string) bool {
	var cart model.CartItem
	dao.DB.Table("cart_items").Where("nick_name=?", nickname).First(&cart)
	if cart.NickName == nickname {
		return true
	}
	return false
}
