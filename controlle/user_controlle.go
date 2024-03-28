package controlle

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go_shop/bcrypt"
	"go_shop/dao"
	"go_shop/model"
	"go_shop/util"
	"strconv"
)

// UserLogin 登录逻辑函数
func UserLogin(userID, password string) (string, error) {
	// 验证用户信息是否符合要求
	if ok, err := LoginMsgIsOk(model.User{UserId: userID, PassWord: password}); !ok {
		return "", fmt.Errorf("用户信息不符合要求: %v", err)
	}

	var user model.User
	// 根据userID查询用户信息
	err := dao.DB.Where("sid=?", userID).First(&user).Error
	if err != nil {
		return "", fmt.Errorf("用户不存在：%v", err)
	}
	if user.PassWord == "123456" {
		token, err := util.GenToken(user.UserId)
		if err != nil {
			return "", fmt.Errorf("生成token失败：%v", err)
		}
		return token, nil
	}
	// 验证密码是否正确
	if !bcrypt.ComparePasswords(password, user.PassWord) {
		return "", errors.New("密码不正确")
	}

	// 生成token
	token, err := util.GenToken(user.UserId)
	if err != nil {
		return "", fmt.Errorf("生成token失败：%v", err)
	}

	logrus.Info("用户登录成功")
	return token, nil
}

// LoginMsgIsOk 判断用户信息是否符合要求
func LoginMsgIsOk(user model.User) (bool, error) {
	if len(user.UserId) != 10 {
		return false, errors.New("userid的长度不符")
	}

	if len(user.PassWord) < 5 || len(user.PassWord) > 14 {
		return false, errors.New("用户密码应在8和14位")
	}

	return true, nil
}

func Login(c *fiber.Ctx) error {

	type Req struct {
		Userid   string `json:"sid"`
		Password string `json:"password"`
	}

	var (
		req  Req
		user model.User
	)

	err := c.BodyParser(&req)
	if err != nil {
		return util.Resp400(c, "参数验证错误")
	}

	token, err := UserLogin(req.Userid, req.Password)

	if err != nil {
		return util.Resp400(c, fmt.Errorf("登录失败 : %v", err))
	}

	dao.DB.Table("users").
		Where("sid = ?", req.Userid).
		First(&user)

	user.PassWord = ""
	return util.Resp200(c, fiber.Map{
		"msg":   "登陆成功",
		"token": token,
		"info":  user,
	})
}

type blackUser struct {
	UserId string `json:"sid"`
	T      string `json:"time"`
}

func KickUser(c *fiber.Ctx) error {
	var buser blackUser
	err := c.BodyParser(&buser)
	// userid和time 绑定 接受json Param from-data三种参数绑定方式
	if err != nil || buser.UserId == "" {
		buser.UserId = c.Query("sid")
		buser.T = c.Query("time")
		if buser.UserId == "" && buser.T == "" {
			buser.UserId = c.FormValue("sid")
			buser.T = c.FormValue("time")
			if buser.UserId == "" || buser.T == "" {
				return util.Resp400(c, "userid或time参数不规范")
			}
		}
	}
	i, err := strconv.ParseInt(buser.T, 10, 64)
	if err != nil {
		logrus.Info("time不符合规范")
		return util.Resp400(c, "拉黑时间不规范")
	}

	if err := util.KickUser(buser.UserId, i); err != nil {
		return util.Resp400(c, fmt.Errorf("error :%v", err))
	}

	return util.Resp200(c, "将用户加入黑名单")
}

// UpdatePassword 修改密码
func UpdatePassword(c *fiber.Ctx) error {
	type Req struct {
		Value string `json:"value"`
	}
	var (
		err error
		req Req
	)
	userID, st := c.Locals("sid").(string)
	if !st {
		return errors.New("类型断言失败")
	}
	err = c.BodyParser(&req)
	if err != nil {
		return util.Resp400(c, "参数绑定失败")
	}
	newPwd, err := bcrypt.HashPassword(req.Value)
	if err != nil {
		return util.Resp500(c, "加密失败")
	}
	err = dao.DB.Table("users").
		Where("sid = ?", userID).
		UpdateColumn("password", newPwd).
		Error
	if err != nil {
		return util.Resp500(c, "修改密码失败")
	}
	return util.Resp200(c, "修改成功")
}

// Information 个人信息
func Information(c *fiber.Ctx) error {
	userID, st := c.Locals("sid").(string)
	if !st {
		return errors.New("类型断言失败")
	}
	var (
		user model.User
		err  error
	)

	err = dao.DB.Table("users").
		Where("sid=?", userID).
		First(&user).
		Error
	if err != nil {
		return util.Resp400(c, "查询失败")
	}
	return util.Resp200(c, user)
}
