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
		return fmt.Errorf("已添加过购物车")
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

// UpdateItem 商品状态修改
//func UpdateItem(nickname string, userid string, num int16) error {
//	var item model.CartItem
//	tx := dao.DB.Table("cart_items").Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			tx.Rollback()
//		}
//	}()
//	if err := dao.DB.Where("user_id=? and nick_name=?", userid, nickname).First(&item).Error; err != nil {
//		return fmt.Errorf("cart.go 88 没有查询到该商品 err : %v", err)
//	}
//	item. = num
//
//	if err := dao.DB.Save(&item).Error; err != nil {
//		tx.Rollback()
//		return fmt.Errorf("cart.go 93 数据库修改出现错误 err:%v", err)
//	}
//
//	if err := tx.Commit().Error; err != nil {
//		tx.Rollback()
//		return fmt.Errorf("事务提交错误 %v", err)
//	}
//	logrus.Info("购物车商品信息修改成功")
//
//	return nil
//}

// DeleteItem 删除商品
func DeleteItem(nickname string, userid string) error {
	tx := dao.DB.Table("cart_items").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := dao.DB.Where("nike_name=? and user_id=?", nickname, userid).Delete(&model.CartItem{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交错误 %v", err)
	}
	logrus.Info("删除成功")
	return nil
}

// SearchItem 查询userid下购物车中所有商品的信息
func SearchItem(userid string) ([]model.CartItem, error) {
	tx := dao.DB.Table("cart_items").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var items []model.CartItem
	if err := dao.DB.Where("user_id=?", userid).Find(&items).Error; err != nil {
		tx.Rollback()
		return items, fmt.Errorf("查询购物车失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return items, fmt.Errorf("事务提交错误 %v", err)
	}
	logrus.Info("查询成功")
	return items, nil
}
