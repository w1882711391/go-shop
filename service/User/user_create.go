package User

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go_shop/bcrypt"
	"go_shop/dao"
	"go_shop/model"
)

type NewUser struct {
}

// UserCreate 注册逻辑函数
func (n *NewUser) UserCreate(user model.User) (string, error) {
	if ok, err := n.UserMsgIsOk(user); !ok {
		return "", err
	}

	//开启事务 防止并发注册同一id
	tx := dao.DB.Table("users").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var newUser model.User
	//	如果错误不为空的话 代表查到了
	dao.DB.Where("user_id=?", user.UserId).First(&newUser)
	if newUser.UserId == user.UserId {
		tx.Rollback()
		return "", fmt.Errorf("该userid已存在用户")
	}

	//开始注册逻辑
	password, err := bcrypt.HashPassword(user.PassWord)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("密码加密出现错误:%v", err)
	}
	user.PassWord = password

	//创建用户
	if err := dao.DB.Create(&user).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("用户创建失败:%v", err)
	}

	// 生成token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.UserId

	// 签名token
	signingKey := []byte("cui cui") // 替换为你自己的密钥
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("生成token失败:%v", err)
	}
	logrus.Info("用户创建成功")
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("事务提交失败:%v", err)
	}
	return signedToken, nil
}

// UserMsgIsOk 判断用户信息是否符合要求
func (n *NewUser) UserMsgIsOk(user model.User) (bool, error) {
	if len(user.UserId) < 6 || len(user.UserId) > 10 {
		return false, errors.New("userid的长度应大于6小于10")
	}

	if len(user.PassWord) < 8 || len(user.PassWord) > 14 {
		return false, errors.New("用户密码应在8和14位")
	}

	if len(user.PhoneNumber) != 11 {
		return false, errors.New("您的手机号信息不符合规范")
	}
	return true, nil
}
