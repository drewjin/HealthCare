package models

import (
	"gorm.io/gorm"
)

type Family struct {
	gorm.Model
	UserID       uint   `gorm:"not null;uniqueIndex:idx_user_relative;column:user_id"`
	RelativeID   uint   `gorm:"not null;uniqueIndex:idx_user_relative;column:relative_id"`
	Relationship string `gorm:"type:varchar(50);not null;column:relationship"`
	Status       uint8  `gorm:"type:tinyint(1);not null;default:0;column:status"` // 0: pending, 1: accepted

	// Relations
	ThisUser User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Relative User `gorm:"foreignKey:RelativeID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
