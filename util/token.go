package util

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Claims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 这里额外记录sid，所以要自定义结构体
type Claims struct {
	UserId string `json:"sid"`
	jwt.RegisteredClaims
}

// mySecret 密钥
var mySecret = []byte("cuicui")

// GenToken 生成 Token
func GenToken(userid string) (string, error) {
	c := Claims{
		UserId: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "cuicui",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3000 * time.Minute)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
		}}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenStr, err := token.SignedString(mySecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func GetPayload(token string) (string, error) {
	parser := jwt.NewParser()
	var claims Claims
	_, _, err := parser.ParseUnverified(token, &claims)
	return claims.UserId, err
}

func VerifyToken(token string) error {
	accessJwtKey := []byte("cuicui")
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return accessJwtKey, nil
	})
	return err
}
