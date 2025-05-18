package models

import "gorm.io/gorm"

// 档案信息表，存储user，organization，plan,item，item_values的信息，显示视图另外显示
/*
	id	user_id	plan_id	item_id	item_value
	1	1		1		1		正常
	2	1		1		2		正常
	3	2		1		1		正常
	4	2		2		2		正常
	5	2		1		3		正常
*/
type UserHealthItem struct {
	gorm.Model
	RelationUserId       uint   `gorm:"not null;index;column:user_id"`
	RelationPlanId       uint   `gorm:"not null;index;column:plan_id"`
	RelationHealthItemId uint   `gorm:"not null;index;column:health_item_id"`
	ItemValue            string `gorm:"type:varchar(100);not null;index;column:item_value"`

	// Relations
	ThisUser      User       `gorm:"foreignKey:RelationUserId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ThisPlan      Plan       `gorm:"foreignKey:RelationPlanId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ThisHeathItem HealthItem `gorm:"foreignKey:RelationHealthItemId;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
