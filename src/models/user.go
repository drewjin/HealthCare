package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique;->;<-:create;column:username"`
	Password string `gorm:"type:varchar(255);not null;->;<-:create;column:password"`
	Name     string `gorm:"type:varchar(50);not null;<-:create;column:name"`
	Gender   string `gorm:"type:char(1);not null;->;<-:create;column:gender"`
	Birthday string `gorm:"type:date;not null;->;<-:create;column:birthday"`
	Phone    string `gorm:"type:varchar(11);not null;->;<-:create;column:phone"`
	Email    string `gorm:"type:varchar(100);->;<-:create;column:email"`
	Address  string `gorm:"type:varchar(255);->;<-:create;column:address"`
}
