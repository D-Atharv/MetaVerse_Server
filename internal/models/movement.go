package models

import "gorm.io/gorm"

type Position struct {
	gorm.Model
	UserID string  `json:"user_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;references:ID"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}
