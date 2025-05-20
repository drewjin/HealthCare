package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserPackagesByInstitution 获取机构的用户套餐列表
func GetUserPackagesByInstitution(c *gin.Context) {
	// 验证用户是否是机构用户
	username := c.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// 检查是否为机构用户 (UserType = 3)
	if user.UserType != 3 && user.UserType != 2 { // 不是机构用户也不是管理员
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: only institution users can view user packages",
		})
		return
	}

	// 获取机构ID
	var institutionID uint
	if user.UserType == 3 { // 如果是机构用户
		var institution models.Institution
		if err := global.DB.Where("user_id = ?", user.ID).Select("id").First(&institution).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Institution not found for this user",
			})
			return
		}
		institutionID = institution.ID
	} else {
		// 如果是管理员，可以指定机构ID
		institutionIDStr := c.Query("institution_id")
		if institutionIDStr != "" {
			id, err := strconv.ParseUint(institutionIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid institution ID",
				})
				return
			}
			institutionID = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Institution ID is required for admin users",
			})
			return
		}
	}

	// 构造查询参数
	query := global.DB.Table("user_packages").
		Select(`user_packages.id, 
                 user_packages.user_id, 
                 user_packages.plan_id, 
                 user_packages.institution_id, 
                 user_packages.status,
                 users.name as user_name, 
                 plans.plan_name`).
		Joins("JOIN users ON user_packages.user_id = users.id").
		Joins("JOIN plans ON user_packages.plan_id = plans.id").
		Where("user_packages.institution_id = ?", institutionID)

	// 处理搜索参数
	if userName := c.Query("user_name"); userName != "" {
		query = query.Where("users.name LIKE ?", "%"+userName+"%")
	}

	if planName := c.Query("plan_name"); planName != "" {
		query = query.Where("plans.plan_name LIKE ?", "%"+planName+"%")
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status, err := strconv.ParseUint(statusStr, 10, 8)
		if err == nil {
			query = query.Where("user_packages.status = ?", status)
		}
	}

	// 定义结果结构
	type UserPackageInfo struct {
		ID             uint   `json:"id"`
		UserID         uint   `json:"user_id"`
		PlanID         uint   `json:"plan_id"`
		InstitutionID  uint   `json:"institution_id"`
		Status         uint8  `json:"status"`
		UserName       string `json:"user_name"`
		PlanName       string `json:"plan_name"`
		CompletedItems int    `json:"completed_items"`
		TotalItems     int    `json:"total_items"`
	}

	var userPackages []UserPackageInfo

	// 执行查询
	if err := query.Find(&userPackages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user packages: " + err.Error(),
		})
		return
	}

	// 补充每个套餐的完成项目数量
	for i := range userPackages {
		// 获取套餐总项目数
		var totalItems int64
		if err := global.DB.Model(&models.PlanHeathItem{}).
			Where("plan_id = ?", userPackages[i].PlanID).
			Count(&totalItems).Error; err != nil {
			continue
		}
		userPackages[i].TotalItems = int(totalItems)

		// 获取用户已完成的项目数
		var completedItems int64
		if err := global.DB.Model(&models.UserHealthItem{}).
			Where("user_id = ? AND plan_id = ? AND item_value != ''", 
				userPackages[i].UserID, userPackages[i].PlanID).
			Count(&completedItems).Error; err != nil {
			continue
		}
		userPackages[i].CompletedItems = int(completedItems)
	}

	c.JSON(http.StatusOK, gin.H{
		"packages": userPackages,
	})
}
