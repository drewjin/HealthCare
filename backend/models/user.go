package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique;index;column:username"`
	Password string `gorm:"type:varchar(255);not null;column:password"`
	Name     string `gorm:"type:varchar(50);not null;column:name"`
	Gender   string `gorm:"type:char(1);not null;->;column:gender"`
	Birthday string `gorm:"type:date;not null;column:birthday"`
	Phone    string `gorm:"type:varchar(11);not null;column:phone"`
	Email    string `gorm:"type:varchar(100);column:email"`
	Address  string `gorm:"type:varchar(255);column:address"`
}

type RolePermission struct {
	gorm.Model
	UserID     uint `gorm:"not null;index;column:user_id"`
	Permission uint `gorm:"not null;column:permission"`

	// Relations
	ThisUser User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

type Family struct {
	gorm.Model
	UserID       uint    `gorm:"not null;uniqueIndex:idx_user_relative;column:user_id"`
	RelativeID   uint    `gorm:"not null;uniqueIndex:idx_user_relative;column:relative_id"`
	Relationship string  `gorm:"type:varchar(50);not null;column:relationship"`
	Status       uintptr `gorm:"type:tinyint(1);not null;column:status"`

	// Relations
	ThisUser User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Relative User `gorm:"foreignKey:RelativeID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
