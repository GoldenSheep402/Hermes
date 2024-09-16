package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

var secret = []byte("todo")

type Info struct {
	UID            string
	OrgID          string
	IsRefreshToken bool
}

type JWTClaims struct {
	Info Info
	jwt.StandardClaims
}

const (
	AccessTokenExpireIn  = time.Hour * 24
	RefreshTokenExpireIn = time.Hour * 24 * 30
)

// GenToken 生成JWT
func GenToken(info Info, expire ...time.Duration) (token string, err error) {
	if string(secret) == ("todo") {
		panic("JWT secret didn't change")
	}
	if len(expire) == 0 {
		expire = append(expire, AccessTokenExpireIn)
	}
	c := JWTClaims{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire[0]).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "jframe",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(secret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
