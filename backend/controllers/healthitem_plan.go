package controllers

import (
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPlanHealthItems 获取套餐关联的健康项目
func GetPlanHealthItems(c *gin.Context) {
	planID := c.Param("plan_id")
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
		ID              uint   `json:"id"`
		HealthItemID    uint   `json:"health_item_id"`
		ItemName        string `json:"item_name"`
		ItemDescription string `json:"item_description"`
	}

	err = global.DB.Model(&models.PlanHeathItem{}).
		Select("plan_heath_items.id, plan_heath_items.health_item_id, health_items.item_name, plan_heath_items.item_description").
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

// AssociatePlanHealthItems 将健康项目关联到套餐
func AssociatePlanHealthItems(c *gin.Context) {
	planID := c.Param("plan_id")
	pid, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的套餐ID",
		})
		return
	}

	// 验证套餐存在
	var plan models.Plan
	if err := global.DB.First(&plan, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在",
		})
		return
	}

	// 验证用户是否有权限操作该套餐（机构用户只能操作自己的套餐）
	username := c.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 管理员可以操作任何套餐，机构用户只能操作自己的套餐
	if user.UserType == 3 { // 如果是机构用户
		var institution models.Institution
		if err := global.DB.Where("user_id = ?", user.ID).First(&institution).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "您没有关联的机构",
			})
			return
		}

		if institution.ID != plan.RelationInstitutionID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "您没有权限操作此套餐",
			})
			return
		}
	} else if user.UserType != 2 && user.UserType != 1 { // 不是管理员或超级管理员
		c.JSON(http.StatusForbidden, gin.H{
			"error": "您没有权限进行此操作",
		})
		return
	}

	// 解析请求体
	type HealthItemInput struct {
		HealthItemID    uint   `json:"health_item_id"`
		ItemDescription string `json:"item_description"`
	}

	var input []HealthItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 开始事务
	tx := global.DB.Begin()

	// 先删除现有关联
	if err := tx.Where("plan_id = ?", pid).Delete(&models.PlanHeathItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除现有关联失败: " + err.Error(),
		})
		return
	}

	// 创建新关联
	for _, item := range input {
		// 验证健康项目存在
		var healthItem models.HealthItem
		if err := tx.First(&healthItem, item.HealthItemID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("ID为%d的健康项目不存在", item.HealthItemID),
			})
			return
		}

		// 创建关联
		planHealthItem := models.PlanHeathItem{
			RelationPlanId:       uint(pid),
			RelationHealthItemId: item.HealthItemID,
			ItemDescription:      item.ItemDescription,
		}

		if err := tx.Create(&planHealthItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "创建健康项目关联失败: " + err.Error(),
			})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "提交事务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "套餐关联健康项目成功",
		"count":   len(input),
	})
}

// AddHealthItemToPlan 添加健康检查项目到套餐中
func AddHealthItemToPlan(c *gin.Context) {
	planID := c.Param("plan_id")
	pid, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的套餐ID",
		})
		return
	}

	// 验证套餐存在
	var plan models.Plan
	if err := global.DB.First(&plan, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在",
		})
		return
	}

	// 验证用户是否有权限操作该套餐（机构用户只能操作自己的套餐）
	username := c.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 管理员可以操作任何套餐，机构用户只能操作自己的套餐
	if user.UserType == 3 { // 如果是机构用户
		var institution models.Institution
		if err := global.DB.Where("user_id = ?", user.ID).First(&institution).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "您没有关联的机构",
			})
			return
		}

		if institution.ID != plan.RelationInstitutionID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "您没有权限操作此套餐",
			})
			return
		}
	} else if user.UserType != 2 && user.UserType != 1 { // 不是管理员或超级管理员
		c.JSON(http.StatusForbidden, gin.H{
			"error": "您没有权限进行此操作",
		})
		return
	}

	// 解析请求体
	var input struct {
		HealthItemID    uint   `json:"health_item_id" binding:"required"`
		ItemDescription string `json:"item_description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 验证健康项目存在
	var healthItem models.HealthItem
	if err := global.DB.First(&healthItem, input.HealthItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "健康项目不存在",
		})
		return
	}

	// 检查项目是否已经关联到套餐
	var existingItem models.PlanHeathItem
	result := global.DB.Where("plan_id = ? AND health_item_id = ?", pid, input.HealthItemID).First(&existingItem)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "该健康项目已关联到此套餐",
		})
		return
	}

	// 创建新关联
	planHealthItem := models.PlanHeathItem{
		RelationPlanId:       uint(pid),
		RelationHealthItemId: input.HealthItemID,
		ItemDescription:      input.ItemDescription,
		ItemMetrics:          "", // 初始为空
	}

	if err := global.DB.Create(&planHealthItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建健康项目关联失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功添加健康项目到套餐",
		"id":      planHealthItem.ID,
	})
}

// GetUserPlanHealthItems 获取用户特定套餐的健康数据
func GetUserPlanHealthItems(c *gin.Context) {
	userID := c.Param("user_id")
	planID := c.Param("plan_id")

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

	// 验证用户有权限查看此数据
	username := c.GetString("username")
	var currentUser models.User
	if err := global.DB.Where("username = ?", username).First(&currentUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权",
		})
		return
	}

	// 处理权限：管理员可以查看任何用户，普通用户只能查看自己，机构用户可以查看选择了其套餐的用户
	var hasPermission bool
	if currentUser.UserType == 1 || currentUser.UserType == 2 { // 超级管理员或管理员
		hasPermission = true
	} else if currentUser.ID == uint(uid) { // 查看自己的数据
		hasPermission = true
	} else if currentUser.UserType == 3 { // 机构用户
		// 查询该机构是否拥有该套餐，以及该用户是否选择了该套餐
		var count int64
		global.DB.Model(&models.Plan{}).
			Joins("JOIN institutions ON plans.institution_id = institutions.id").
			Joins("JOIN user_packages ON plans.id = user_packages.plan_id").
			Where("institutions.user_id = ? AND plans.id = ? AND user_packages.user_id = ?", currentUser.ID, pid, uid).
			Count(&count)

		hasPermission = count > 0
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "您没有权限查看此数据",
		})
		return
	}

	// 查询用户该套餐下的健康项目数据
	type HealthItemData struct {
		ItemID      uint   `json:"item_id"`
		ItemName    string `json:"item_name"`
		ItemValue   string `json:"item_value"`
		Description string `json:"description"`
	}

	var healthData []HealthItemData

	err = global.DB.Model(&models.UserHealthItem{}).
		Select("user_health_items.health_item_id as item_id, health_items.item_name, user_health_items.item_value, plan_heath_items.item_description as description").
		Joins("JOIN health_items ON user_health_items.health_item_id = health_items.id").
		Joins("LEFT JOIN plan_heath_items ON plan_heath_items.plan_id = user_health_items.plan_id AND plan_heath_items.health_item_id = user_health_items.health_item_id").
		Where("user_health_items.user_id = ? AND user_health_items.plan_id = ?", uid, pid).
		Find(&healthData).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取健康数据失败: " + err.Error(),
		})
		return
	}

	// 获取套餐和用户信息
	var plan models.Plan
	if err := global.DB.First(&plan, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "套餐不存在",
		})
		return
	}

	var targetUser models.User
	if err := global.DB.Select("id, username, name, gender, birthday").First(&targetUser, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       targetUser.ID,
			"username": targetUser.Username,
			"name":     targetUser.Name,
			"gender":   targetUser.Gender,
			"birthday": targetUser.Birthday,
		},
		"plan": gin.H{
			"id":   plan.ID,
			"name": plan.PlanName,
		},
		"health_data": healthData,
	})
}

// GetUserInstitutionHealthItems 按机构获取用户的健康数据
func GetUserInstitutionHealthItems(c *gin.Context) {
	userID := c.Param("user_id")
	institutionID := c.Param("institution_id")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	iid, err := strconv.ParseUint(institutionID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的机构ID",
		})
		return
	}

	// 验证用户有权限查看此数据
	username := c.GetString("username")
	var currentUser models.User
	if err := global.DB.Where("username = ?", username).First(&currentUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权",
		})
		return
	}

	// 处理权限：管理员可以查看任何用户，普通用户只能查看自己，机构用户可以查看选择了其套餐的用户
	var hasPermission bool
	if currentUser.UserType == 1 || currentUser.UserType == 2 { // 超级管理员或管理员
		hasPermission = true
	} else if currentUser.ID == uint(uid) { // 查看自己的数据
		hasPermission = true
	} else if currentUser.UserType == 3 { // 机构用户
		// 查询该用户是否为该机构的管理员
		var institution models.Institution
		global.DB.Where("id = ? AND user_id = ?", iid, currentUser.ID).First(&institution)

		hasPermission = institution.ID > 0
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "您没有权限查看此数据",
		})
		return
	}

	// 获取该机构下用户选择的所有套餐
	var userPackages []struct {
		PlanID   uint   `json:"plan_id"`
		PlanName string `json:"plan_name"`
	}

	err = global.DB.Model(&models.UserPackage{}).
		Select("user_packages.plan_id, plans.plan_name").
		Joins("JOIN plans ON user_packages.plan_id = plans.id").
		Where("user_packages.user_id = ? AND plans.institution_id = ?", uid, iid).
		Find(&userPackages).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户套餐失败: " + err.Error(),
		})
		return
	}

	if len(userPackages) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户未选择该机构的套餐",
		})
		return
	}

	// 获取机构信息
	var institution models.Institution
	if err := global.DB.First(&institution, iid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "机构不存在",
		})
		return
	}

	// 获取用户信息
	var targetUser models.User
	if err := global.DB.Select("id, username, name, gender, birthday").First(&targetUser, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 获取所有套餐的健康数据
	var planHealthData []struct {
		PlanID     uint                     `json:"plan_id"`
		PlanName   string                   `json:"plan_name"`
		HealthData []map[string]interface{} `json:"health_data"`
	}

	for _, pkg := range userPackages {
		var healthItems []struct {
			ItemID    uint   `json:"item_id"`
			ItemName  string `json:"item_name"`
			ItemValue string `json:"item_value"`
		}

		global.DB.Model(&models.UserHealthItem{}).
			Select("user_health_items.health_item_id as item_id, health_items.item_name, user_health_items.item_value").
			Joins("JOIN health_items ON user_health_items.health_item_id = health_items.id").
			Where("user_health_items.user_id = ? AND user_health_items.plan_id = ?", uid, pkg.PlanID).
			Find(&healthItems)

		var healthData []map[string]interface{}
		for _, item := range healthItems {
			healthData = append(healthData, map[string]interface{}{
				"item_id":    item.ItemID,
				"item_name":  item.ItemName,
				"item_value": item.ItemValue,
			})
		}

		planHealthData = append(planHealthData, struct {
			PlanID     uint                     `json:"plan_id"`
			PlanName   string                   `json:"plan_name"`
			HealthData []map[string]interface{} `json:"health_data"`
		}{
			PlanID:     pkg.PlanID,
			PlanName:   pkg.PlanName,
			HealthData: healthData,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       targetUser.ID,
			"username": targetUser.Username,
			"name":     targetUser.Name,
			"gender":   targetUser.Gender,
			"birthday": targetUser.Birthday,
		},
		"institution": gin.H{
			"id":   institution.ID,
			"name": institution.InstitutionName,
		},
		"plans": planHealthData,
	})
}
