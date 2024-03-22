package User

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_shop/bcrypt"
	"go_shop/dao"
	"go_shop/model"
	"go_shop/util"
)

// UserLogin 登录逻辑函数
func UserLogin(userID, password string) (string, error) {
	// 验证用户信息是否符合要求
	if ok, err := LoginMsgIsOk(model.User{UserId: userID, PassWord: password}); !ok {
		return "", fmt.Errorf("用户信息不符合要求: %v", err)
	}

	var user model.User
	// 根据userID查询用户信息
	err := dao.DB.Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return "", fmt.Errorf("用户不存在：%v", err)
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

	if len(user.PassWord) < 8 || len(user.PassWord) > 14 {
		return false, errors.New("用户密码应在8和14位")
	}

	return true, nil
}
