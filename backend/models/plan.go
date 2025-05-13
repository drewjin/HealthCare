package models

import "gorm.io/gorm"

// 套餐表，只存储套餐信息，plan_id，plan_name,organization_id(外键)
type Plan struct {
	gorm.Model
	PlanName                string `gorm:"type:varchar(100);not null;unique;index;column:plan_name"`
	RelattionOrganizationID uint   `gorm:"not null;index;column:organization_id"`

	// Relations
	ThisOrganization Organization `gorm:"foreignKey:RelattionOrganizationID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// 体检套餐指标表，存储具体套餐的指标信息，item_id，item_name(主键),只保存信息，便于添加新的指标内容
type HeathItem struct {
	gorm.Model
	ItemName string `gorm:"type:varchar(100);not null;index;column:item_name"`
}

// plan-item对应表, 展示套餐信息
type PlanHeathItem struct {
	gorm.Model
	RelationPlanId      uint `gorm:"not null;index;column:plan_id"`
	RelationHeathItemId uint `gorm:"not null;index;column:heath_item_id"`

	// Relations
	ThisPlan      Plan      `gorm:"foreignKey:RelationPlanId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ThisHeathItem HeathItem `gorm:"foreignKey:RelationHeathItemId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
