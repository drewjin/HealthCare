package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateUserPackageStatus 更新用户套餐状态
func UpdateUserPackageStatus(c *gin.Context) {
	// 获取路径参数
	userID := c.Param("user_id")
	planID := c.Param("plan_id")
	
	// 解析用户ID和套餐ID
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	pid, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的套餐ID",
		})
		return
	}

	// 验证用户是否有权限更新
	// 只有机构用户(UserType=3)或管理员(UserType=2)可以更新
	username := c.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	if user.UserType != 2 && user.UserType != 3 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: only institution users or admins can update package status",
		})
		return
	}

	// 解析请求体中的状态值
	var input struct {
		Status uint8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: status is required",
		})
		return
	}

	// 验证状态值
	if input.Status > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status value: must be 0 (pending) or 1 (completed)",
		})
		return
	}

	// 更新用户套餐状态
	result := global.DB.Model(&models.UserPackage{}).
		Where("user_id = ? AND plan_id = ?", uid, pid).
		Update("status", input.Status)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update package status: " + result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Package not found for this user and plan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Package status updated successfully",
	})
}
