package controlle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_shop/model"
	"go_shop/service/User"
	"go_shop/util"
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

func Login(ctx *fiber.Ctx) error {
	userid := ctx.FormValue("user_id")
	password := ctx.FormValue("pass_word")

	token, err := User.UserLogin(userid, password)

	if err != nil {
		return util.Resp500(ctx, fmt.Errorf("登录失败 : %v", err))
	}

	return util.Resp200(ctx, fiber.Map{
		"msg":   "登陆成功",
		"token": token,
	})
}
