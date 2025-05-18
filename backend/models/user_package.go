package models

import (
	"gorm.io/gorm"
)

// UserPackage represents a package selected by a user
type UserPackage struct {
	gorm.Model
	UserID        uint  `json:"user_id" gorm:"not null;index;column:user_id"`
	InstitutionID uint  `json:"institution_id" gorm:"not null;index;column:institution_id"`
	PlanID        uint  `json:"plan_id" gorm:"not null;index;column:plan_id"`
	Status        uint8 `json:"status" gorm:"type:tinyint(1);not null;default:0;column:status"` // 0: pending, 1: completed

	// Relations
	User        User        `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Institution Institution `json:"institution" gorm:"foreignKey:InstitutionID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Plan        Plan        `json:"plan" gorm:"foreignKey:PlanID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
