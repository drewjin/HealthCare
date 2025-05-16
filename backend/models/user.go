package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique;index;column:username"`
	Password string `gorm:"type:varchar(255);not null;column:password"`
	Name     string `gorm:"type:varchar(50);not null;column:name"`
	Gender   string `gorm:"type:enum('男','女','未知');not null;column:gender"`
	Birthday string `gorm:"type:date;column:birthday"`
	Phone    string `gorm:"type:varchar(11);not null;column:phone"`
	Email    string `gorm:"type:varchar(100);column:email"`
	Address  string `gorm:"type:varchar(255);column:address"`
	UserType uint8  `gorm:"type:tinyint(1);not null;default:1;column:user_type" json:"user_type"` // 1: normal, 2: admin, 3: institution
}

type RolePermission struct {
	gorm.Model
	UserID     uint `gorm:"not null;index;column:user_id"`
	Permission uint `gorm:"not null;column:permission"`

	// Relations
	ThisUser User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
