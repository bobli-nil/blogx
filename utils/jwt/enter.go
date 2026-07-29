package jwt

import (
	"blogx_server/global"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	UserID   uint   `json:"userID"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, userName string) (string, error) {
	expireTime := time.Now().Add(time.Duration(global.Conf.Jwt.Expire) * time.Second)

	claims := MyClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Conf.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Conf.Jwt.Secret))
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(global.Conf.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseTokenByGin 从gin中获取token
func ParseTokenByGin(c *gin.Context) (*MyClaims, error) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	return ParseToken(token)
}
