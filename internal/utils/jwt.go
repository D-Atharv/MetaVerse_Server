package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"os"
)

func GenerateJWT(username string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
