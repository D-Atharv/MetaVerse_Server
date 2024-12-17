package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}
