package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RecoverDeletedPlan 恢复已删除的套餐（如果应用了软删除）
func RecoverDeletedPlan(c *gin.Context) {
	// 获取套餐ID
	planID := c.Param("id")
	pid, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的套餐ID",
		})
		return
	}

	// 验证用户是否有权限恢复套餐（仅限管理员）
	username := c.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// 只有管理员可以恢复已删除的套餐
	if user.UserType != 2 { // 2表示管理员
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: only admins can recover deleted plans",
		})
		return
	}

	// 检查套餐是否存在（包括已删除的）
	var plan models.Plan
	if err := global.DB.Unscoped().First(&plan, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在，无法恢复",
		})
		return
	}

	// 检查套餐是否已被删除
	if plan.DeletedAt.Time.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "套餐未被删除，无需恢复",
		})
		return
	}

	// 恢复套餐
	if err := global.DB.Unscoped().Model(&models.Plan{}).Where("id = ?", pid).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "恢复套餐失败: " + err.Error(),
		})
		return
	}

	// 恢复相关的PlanHealthItem（如果需要）
	if err := global.DB.Unscoped().Model(&models.PlanHeathItem{}).Where("plan_id = ?", pid).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "恢复套餐项目失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "套餐恢复成功",
		"plan_id": pid,
	})
}

// GetDeletedPlans 获取已删除的套餐列表
func GetDeletedPlans(c *gin.Context) {
	// 仅限管理员访问
	username := c.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// 只有管理员可以查看已删除的套餐
	if user.UserType != 2 { // 2表示管理员
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: only admins can view deleted plans",
		})
		return
	}

	// 获取已删除的套餐列表
	var plans []struct {
		ID            uint    `json:"id"`
		PlanName      string  `json:"plan_name"`
		InstitutionID uint    `json:"institution_id"`
		InstName      string  `json:"institution_name"`
		PlanPrice     float64 `json:"plan_price"`
		DeletedAt     string  `json:"deleted_at"`
		ItemCount     int64   `json:"item_count"`
	}

	err := global.DB.Table("plans").
		Select("plans.id, plans.plan_name, plans.institution_id, institutions.institution_name as inst_name, plans.plan_price, plans.deleted_at").
		Joins("JOIN institutions ON plans.institution_id = institutions.id").
		Where("plans.deleted_at IS NOT NULL").
		Find(&plans).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取已删除套餐列表失败: " + err.Error(),
		})
		return
	}

	// 获取每个套餐的项目数量
	for i := range plans {
		var count int64
		global.DB.Unscoped().Model(&models.PlanHeathItem{}).
			Where("plan_id = ?", plans[i].ID).
			Count(&count)
		plans[i].ItemCount = count
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted_plans": plans,
	})
}
