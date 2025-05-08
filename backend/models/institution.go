package models

import (
	"gorm.io/gorm"
)

type Institution struct {
	gorm.Model
	InstitutionName          string `json:"institution_name" gorm:"type:varchar(50);not null;unique;column:institution_name"`
	InstitutionAddress       string `json:"institution_address" gorm:"type:varchar(100);not null;column:institution_address"`
	InstitutionQualification string `json:"institution_qualification" gorm:"type:varchar(50);column:institution_qualification"`
	ExaminationPackage       string `json:"examination_package" gorm:"type:varchar(255);column:examination_package"`
	UserID                   uint   `json:"user_id" gorm:"not null;index;column:user_id"`
	Status                   uint8  `json:"status" gorm:"type:tinyint(1);not null;default:0;column:status"` // 0: pending, 1: approved, 2: rejected

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
