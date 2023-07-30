package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id int, fullname string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["fullname"] = fullname
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}
