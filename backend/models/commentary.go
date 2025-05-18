package models

import "gorm.io/gorm"

type Commentary struct {
	gorm.Model
	RelationUserId uint   `gorm:"not null;index;column:user_id"`
	RelationPlanId uint   `gorm:"not null;index;column:plan_id"`
	Commentary     string `gorm:"type:varchar(255);not null;column:commentary"`

	// Relations
	ThisUser User `gorm:"foreignKey:RelationUserId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ThisPlan Plan `gorm:"foreignKey:RelationPlanId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
