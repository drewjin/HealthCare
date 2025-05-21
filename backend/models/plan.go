package models

import "gorm.io/gorm"

// 套餐表，只存储套餐信息，plan_id，plan_name,organization_id(外键)以及价格、描述等
type Plan struct {
	gorm.Model
	PlanName              string  `gorm:"type:varchar(100);not null;unique;index;column:plan_name" json:"plan_name"`
	RelationInstitutionID uint    `gorm:"not null;index;column:institution_id" json:"institution_id"`
	PlanPrice             float64 `gorm:"type:decimal(10,2);default:0;column:plan_price" json:"plan_price"`
	Description           string  `gorm:"type:text;column:description" json:"description"`
	SuitableFor           string  `gorm:"type:varchar(255);column:suitable_for" json:"suitable_for"`

	// Relations
	ThisInstitution Institution `gorm:"foreignKey:RelationInstitutionID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;" json:"-"`
}

// 体检套餐指标表，存储具体套餐的指标信息，item_id，item_name(主键),只保存信息，便于添加新的指标内容
type HealthItem struct {
	gorm.Model
	ItemName       string `gorm:"type:varchar(512);index;column:item_name"`
	UserID         uint   `gorm:"index;column:user_id"`
	UserHealthInfo string `gorm:"type:varchar(512);index;column:user_health_info"`
}

// plan-item对应表, 展示套餐信息
type PlanHeathItem struct {
	gorm.Model
	RelationPlanId       uint   `gorm:"not null;index;column:plan_id"`
	RelationHealthItemId uint   `gorm:"not null;index;column:health_item_id"`
	ItemDescription      string `gorm:"type:varchar(100);index;column:item_description"`
	ItemMetrics          string `gorm:"type:text;column:item_metrics"` // Stores key-value metrics as JSON

	// Relations
	ThisPlan      Plan       `gorm:"foreignKey:RelationPlanId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ThisHeathItem HealthItem `gorm:"foreignKey:RelationHealthItemId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
