package models

import "gorm.io/gorm"

/*
	机构id	必填
	机构名称 必填
	机构地址 必填
	机构电话 必填
*/
type Organization struct {
	gorm.Model
	OrganizationName    string `gorm:"type:varchar(100);not null;unique;index;column:organization_name"`
	OrganizationAddress string `gorm:"type:varchar(255);not null;column:organization_address"`
	OrganizationPhone   string `gorm:"type:varchar(11);not null;column:organization_phone"`
}
