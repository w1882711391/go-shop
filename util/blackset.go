package util

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"time"
)

func IsKick() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user model.User
		user.UserId = c.FormValue("user_id")
		if user.UserId == "" {
			if err := c.BodyParser(&user); err != nil {
				return Resp401(c, fmt.Sprintf("参数存在问题 : %v", err))
			}
		}
		exists, err := dao.Client.Exists("blacklist:" + user.UserId).Result()
		if err != nil {
			logrus.Info("blackset 21 error : ", err)
			return Resp500(c, fmt.Errorf("服务器错误 查询黑名单失败： %v", err))
		}
		if exists == 1 {
			logrus.Info("用户在黑名单中 请求拒绝")
			return Resp403(c, "用户在黑名单中 请求拒绝")
		}

		return c.Next()
	}
}

// KickUser 踢出用户
func KickUser(userid string, t int64) error {
	var user model.User
	if err := dao.DB.Table("users").Where("user_id=?", userid).First(&user).Error; err != nil {
		return errors.New("没有该用户")
	}
	err := dao.Client.Set("blacklist:"+userid, "true", time.Duration(t)*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("userid加入黑名单失败 %v", err)
	}
	logrus.Info("将用户加入黑名单成功")
	return nil
}
