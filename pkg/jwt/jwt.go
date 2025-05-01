package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(viper.GetString("jwt.secret"))

type MyClaims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

func GenToken(ID int64) (string, error) {
	c := MyClaims{
		ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 86400 * time.Duration(viper.GetInt("jwt.expire_day"))).Unix(),
			Issuer:    "App",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid jwt token")
}
