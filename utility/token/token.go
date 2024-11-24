package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	Username             string `json:"username"`
	jwt.RegisteredClaims        // v5版本新加的方法
}

var SigningKey = []byte("whatthefuck123weishenmebuneng123") // 密钥必须要长，要达到这个位数，26个英文字母是不行的。！！！！

func Token(username string, durationHours int) (string, error) {
	if durationHours <= 0 {
		return "", fmt.Errorf("duration must be greater than zero")
	}
	var claims = JWTToken{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(durationHours) * time.Hour)), // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                               // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                               // 签发时间
			Issuer:    "xuran",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 使用HS256算法
	ss, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}
func ParseJwt(tokenstring string) (*JWTToken, error) {
	t, err := jwt.ParseWithClaims(tokenstring, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})

	if claims, ok := t.Claims.(*JWTToken); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func ValidToken(tokenstring string) (bool, error) {
	res, err := ParseJwt(tokenstring)
	if err != nil {
		return false, err
	}
	timeInstan := jwt.NewNumericDate(time.Now()).Time
	p := timeInstan.Before(res.ExpiresAt.Time)
	if p {
		return true, nil
	} else {
		return false, errors.New("error time")
	}
}
