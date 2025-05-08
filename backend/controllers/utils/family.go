package utils

import (
	"errors"
	"healthcare/global"
	"healthcare/models"
	"strconv"

	"gorm.io/gorm"
)

func UnmarshalUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	return uint(i), nil
}

func CreateFamilyRequest(userID, relativeID uint, email, relationship string) error {
	// 检查是否已经存在关系（包括反向关系）
	var existingFamily models.Family
	if err := global.DB.Where(
		"((user_id = ? AND relative_id = ?) OR (user_id = ? AND relative_id = ?)) AND status = ?",
		userID, relativeID,
		relativeID, userID,
		1,
	).First(&existingFamily).Error; err == nil {
		return errors.New("family relationship already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 检查是否已经有待处理的请求
	if err := global.DB.Where(
		"((user_id = ? AND relative_id = ?) OR (user_id = ? AND relative_id = ?)) AND status = ?",
		userID, relativeID,
		relativeID, userID,
		0,
	).First(&existingFamily).Error; err == nil {
		return errors.New("there is already a pending request for this relationship")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

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
