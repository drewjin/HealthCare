package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique;column:username"`
	Password string `gorm:"type:varchar(255);not null;column:password"`
	Name     string `gorm:"type:varchar(50);not null;column:name"`
	Gender   string `gorm:"type:char(1);not null;->;column:gender"`
	Birthday string `gorm:"type:date;not null;column:birthday"`
	Phone    string `gorm:"type:varchar(11);not null;column:phone"`
	Email    string `gorm:"type:varchar(100);column:email"`
	Address  string `gorm:"type:varchar(255);column:address"`
}
