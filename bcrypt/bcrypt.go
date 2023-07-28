package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	// 生成盐值，并用盐值对密码进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 返回哈希后的密码作为字符串
	return string(hashedPassword), nil
}

// ComparePasswords 比较密码和已加密的密码是否匹配
func ComparePasswords(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
