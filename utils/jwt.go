package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte("supersecret_key_change_me")

func GenerateToken(id string, email string) (string, error) {

	claims := jwt.MapClaims{
		"userId": id,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JWT_SECRET)
}
