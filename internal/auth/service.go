package auth

import (
	"golang.org/x/crypto/bcrypt"
	"server/internal/db"
	"server/internal/models"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func CreateUser(username, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user := models.User{Username: username, Password: hashedPassword}
	return database.DB.Create(&user).Error
}

func AuthenticateUser(username, password string) bool {
	var user models.User
	database.DB.Where("username = ?", username).First(&user)
	return VerifyPassword(user.Password, password)
}
