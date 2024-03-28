package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_shop/dao"
	"go_shop/model"
)

func IsMerchant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, st := c.Locals("sid").(string)
		if !st {
			return Resp500(c, "util merchant.go 12: 类型断言失败")
		}

		var user model.User
		if err := dao.DB.Table("users").Where("sid=?", userID).First(&user).Error; err != nil {
			return Resp401(c, fmt.Errorf("util merchant.go 12 err :%v", err))
		}

		//用户不是商户
		if user.IsMerchant == false {
			return Resp403(c, "用户不是商户 不允许对商品进行操作")
		}
		return c.Next()
	}
}
