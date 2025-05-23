package utils

import (
	"HealthCare/backend/global"
)

func UpdateIt(model interface{}, id interface{}, field string, value interface{}) (bool, error) {
	result := global.DB.Model(model).Where("id = ?", id).Update(field, value)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
