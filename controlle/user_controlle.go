package controlle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go_shop/model"
	"go_shop/service/User"
	"go_shop/util"
	"strconv"
)

func Register(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		return util.Resp401(c, fmt.Sprintf("参数存在问题 : %v", err))
	}
	token, err := User.UserCreate(user)
	if err != nil {
		return util.Resp400(c, fmt.Sprintf("注册存在错误: %v", err))
	}

	return util.Resp200(c, fiber.Map{
		"msg":   "注册成功",
		"token": token,
	})
}

func Login(c *fiber.Ctx) error {
	userid := c.FormValue("user_id")
	password := c.FormValue("pass_word")

	token, err := User.UserLogin(userid, password)

	if err != nil {
		return util.Resp500(c, fmt.Errorf("登录失败 : %v", err))
	}

	return util.Resp200(c, fiber.Map{
		"msg":   "登陆成功",
		"token": token,
	})
}

type blackUser struct {
	UserId string `json:"user_id"`
	t      string `json:"time"`
}

func KickUser(c *fiber.Ctx) error {
	var buser blackUser
	err := c.BodyParser(&buser)
	// userid和time 绑定 接受json Param from-data三种参数绑定方式
	if err != nil || buser.UserId == "" {
		buser.UserId = c.Query("user_id")
		buser.t = c.Query("time")
		if buser.UserId == "" && buser.t == "" {
			buser.UserId = c.FormValue("user_id")
			buser.t = c.FormValue("time")
			if buser.UserId == "" || buser.t == "" {
				return util.Resp400(c, "userid或time参数不规范")
			}
		}
	}
	i, err := strconv.ParseInt(buser.t, 10, 64)
	if err != nil {
		logrus.Info("time不符合规范")
		return util.Resp400(c, "拉黑时间不规范")
	}

	if err := util.KickUser(buser.UserId, i); err != nil {
		return util.Resp400(c, fmt.Errorf("error :%v", err))
	}

	return util.Resp200(c, "将用户加入黑名单")
}
