package controlle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_shop/util"
)

func (u *UserHandler) Login(ctx *fiber.Ctx) error {
	userid := ctx.FormValue("user_id")
	password := ctx.FormValue("pass_word")

	token, err := u.User.UserLogin(userid, password)

	if err != nil {
		return util.Resp500(ctx, fmt.Errorf("登录失败 : %v", err))
	}

	return util.Resp200(ctx, fiber.Map{
		"msg":   "登陆成功",
		"token": token,
	})
}
