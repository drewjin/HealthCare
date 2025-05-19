package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserPlans 获取用户的套餐列表
func GetUserPlans(c *gin.Context) {
	userID := c.Param("id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	var userPackages []struct {
		ID              uint   `json:"id"`
		PlanID          uint   `json:"plan_id"`
		PlanName        string `json:"plan_name"`
		InstitutionID   uint   `json:"institution_id"`
		InstitutionName string `json:"institution_name"`
		Status          uint8  `json:"status"`
	}

	err = global.DB.Model(&models.UserPackage{}).
		Select("user_packages.id, user_packages.plan_id, plans.plan_name, user_packages.institution_id, institutions.institution_name, user_packages.status").
		Joins("JOIN plans ON user_packages.plan_id = plans.id").
		Joins("JOIN institutions ON user_packages.institution_id = institutions.id").
		Where("user_packages.user_id = ?", uid).
		Find(&userPackages).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户套餐列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userPackages)
}

// GetPlanItems 获取套餐的体检项目列表
func GetPlanItems(c *gin.Context) {
	planID := c.Param("id")
	pid, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的套餐ID",
		})
		return
	}

	// 先获取套餐名称
	var plan models.Plan
	if err := global.DB.First(&plan, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在",
		})
		return
	}

	// 获取套餐项目
	var planItems []struct {
		ID       uint   `json:"id"`
		ItemName string `json:"item_name"`
	}

	err = global.DB.Model(&models.PlanHeathItem{}).
		Select("plan_heath_items.health_item_id as id, health_items.item_name").
		Joins("JOIN health_items ON plan_heath_items.health_item_id = health_items.id").
		Where("plan_heath_items.plan_id = ?", pid).
		Find(&planItems).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取套餐项目失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plan_id":   pid,
		"plan_name": plan.PlanName,
		"items":     planItems,
	})
}

// GetPlanDetails 获取套餐详情
func GetPlanDetails(c *gin.Context) {
	planID := c.Param("id")
	pid, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的套餐ID",
		})
		return
	}

	var plan models.Plan
	if err := global.DB.First(&plan, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在",
		})
		return
	}

	c.JSON(http.StatusOK, plan)
}
