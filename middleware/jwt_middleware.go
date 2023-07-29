package middleware

import (
	"kajianku_be/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id int, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["fullname"] = name
	claims["exp"] = time.Now().Add(time.Hour * 240).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.SECRET_JWT))
}
