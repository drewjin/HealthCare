package utils

import (
	"errors"
	"healthcare/global"
	"healthcare/models"
	"strconv"
)

func UnmarshalUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	return uint(i), nil
}

func CreateFamilyRequest(userID, relativeID uint, email, relationship string) error {
	family := models.Family{
		UserID:       userID,
		RelativeID:   relativeID,
		Relationship: relationship,
		Status:       0, // 0表示未确认状态
	}

	if err := global.DB.Create(&family).Error; err != nil {
		return err
	}

	return nil
}
