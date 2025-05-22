package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"HealthCare/backend/global"
	"HealthCare/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 查看该用户下的所有健康记录，包含机构和套餐信息
func GetAllItems(ctx *gin.Context) {
	username := ctx.GetString("username")
	var userID uint
	if err := global.DB.Model(&models.User{}).Where("username = ?", username).Select("id").First(&userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 获取用户所选套餐
	var userPackages []models.UserPackage
	if err := global.DB.Preload("Institution").Preload("Plan").Where("user_id = ?", userID).Find(&userPackages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user packages"})
		return
	}

	type record struct {
		PlanID          uint   `json:"plan_id"`
		InstitutionID   uint   `json:"institution_id"`
		InstitutionName string `json:"institution_name"`
		PlanName        string `json:"plan_name"`
		Status          uint8  `json:"status"` // 0: pending, 1: completed
		Items           string `json:"items"`
		ItemCount       int    `json:"item_count"`
		CompletedCount  int    `json:"completed_count"`
	}
	var records []record

	// 对每个套餐聚合健康记录
	for _, pkg := range userPackages {
		// 查询该套餐下的所有健康项目
		var planHealthItems []models.PlanHeathItem
		if err := global.DB.Preload("ThisHeathItem").Where("plan_id = ?", pkg.PlanID).Find(&planHealthItems).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan health items"})
			return
		}

		// 查询该用户该套餐下的所有健康项目记录
		var userHealthItems []models.UserHealthItem
		if err := global.DB.Where("user_id = ? AND plan_id = ?", userID, pkg.PlanID).Find(&userHealthItems).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user health items"})
			return
		}

		// 创建一个映射，方便快速查找用户的检查结果
		userItemValues := make(map[uint]string)
		for _, uhi := range userHealthItems {
			userItemValues[uhi.RelationHealthItemId] = uhi.ItemValue
		}

		// 构建列表
		var list []map[string]string
		completedCount := 0

		for _, phi := range planHealthItems {
			itemValue, exists := userItemValues[phi.RelationHealthItemId]
			if exists && itemValue != "" {
				completedCount++
			}

			list = append(list, map[string]string{
				"item_id":          fmt.Sprintf("%d", phi.RelationHealthItemId),
				"item_name":        phi.ThisHeathItem.ItemName,
				"item_description": phi.ItemDescription,
				"item_value":       itemValue,
			})
		}

		bytes, err := json.Marshal(list)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal items"})
			return
		}

		records = append(records, record{
			PlanID:          pkg.PlanID,
			InstitutionID:   pkg.InstitutionID,
			InstitutionName: pkg.Institution.InstitutionName,
			PlanName:        pkg.Plan.PlanName,
			Status:          pkg.Status,
			Items:           string(bytes),
			ItemCount:       len(planHealthItems),
			CompletedCount:  completedCount,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"records": records})
}

// 查看指定套餐下的体检项目
func GetItemsByPlanID(c *gin.Context) {
	var input struct {
		PlanID uint `form:"plan_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 获取当前用户ID
	username := c.GetString("username")
	var userID uint
	if err := global.DB.Model(&models.User{}).Where("username = ?", username).Select("id").First(&userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 检查是否为机构用户，如果是，获取customerID参数
	var customerID uint
	var isInstitution bool
	var user models.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 如果是机构用户，可以查看所有用户的记录
	isInstitution = user.UserType == 3 // 3表示机构用户

	// 如果是机构用户且提供了customerID参数，使用提供的ID
	if isInstitution && c.Query("customer_id") != "" {
		customerIDStr := c.Query("customer_id")
		customerIDUint, err := strconv.ParseUint(customerIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}
		customerID = uint(customerIDUint)
	} else {
		// 否则使用当前登录用户的ID
		customerID = userID
	}

	type output struct {
		PlanItemID      uint   `json:"plan_item_id"` // PlanHeathItem的ID
		PlanID          uint   `json:"plan_id"`
		PlanName        string `json:"plan_name"`
		ItemID          uint   `json:"item_id"` // HealthItem的ID
		ItemName        string `json:"item_name"`
		ItemDescription string `json:"item_description"`
		ItemValue       string `json:"item_value"`
	}
	var outplanitems []output

	// 先检查该用户是否拥有这个套餐
	var userPackageCount int64
	global.DB.Model(&models.UserPackage{}).
		Where("user_id = ? AND plan_id = ?", customerID, input.PlanID).
		Count(&userPackageCount)

	// 获取套餐详情
	var plan models.Plan
	if err := global.DB.First(&plan, input.PlanID).Error; err != nil {
		// 如果套餐不存在，但我们是机构用户或管理员，仍然尝试获取套餐项目
		if isInstitution || user.UserType == 2 {
			// 检查是否有任何关联到这个套餐的项目
			var planItemCount int64
			if err := global.DB.Model(&models.PlanHeathItem{}).
				Where("plan_id = ?", input.PlanID).
				Count(&planItemCount).Error; err != nil || planItemCount == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found and no items associated"})
				return
			}

			// 创建一个基本的套餐对象供后续使用
			plan = models.Plan{
				Model:    gorm.Model{ID: input.PlanID},
				PlanName: fmt.Sprintf("Plan #%d (Not Found)", input.PlanID),
			}
		} else {
			// 如果是普通用户，并且没有关联到这个套餐，返回404
			if userPackageCount == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found or not associated with this user"})
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
	} else if !isInstitution && user.UserType != 2 && userPackageCount == 0 {
		// 如果套餐存在，但普通用户没有关联到这个套餐，拒绝访问
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this plan"})
		return
	}

	// 查询套餐下的所有体检项目
	var planItems []models.PlanHeathItem
	if err := global.DB.Preload("ThisHeathItem").Where("plan_id = ?", input.PlanID).Find(&planItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan items"})
		return
	}

	// 如果没有找到任何项目，返回空列表
	if len(planItems) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"plan_id":     input.PlanID,
			"plan_name":   plan.PlanName,
			"customer_id": customerID,
			"plan_items":  []output{},
			"message":     "No health items found for this plan",
		})
		return
	}

	// 对于每个套餐项目，获取用户数据(如果存在)
	for _, item := range planItems {
		var itemValue string
		_ = global.DB.Model(&models.UserHealthItem{}).
			Where("user_id = ? AND plan_id = ? AND health_item_id = ?",
				customerID, input.PlanID, item.RelationHealthItemId).
			Select("item_value").
			First(&itemValue).Error

		// 确保item.ThisHeathItem存在
		itemName := "Unknown Item"
		if item.ThisHeathItem.ID > 0 {
			itemName = item.ThisHeathItem.ItemName
		}

		outplanitems = append(outplanitems, output{
			PlanItemID:      item.ID,
			PlanID:          input.PlanID,
			PlanName:        plan.PlanName,
			ItemID:          item.RelationHealthItemId,
			ItemName:        itemName,
			ItemDescription: item.ItemDescription,
			ItemValue:       itemValue,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"plan_id":     input.PlanID,
		"plan_name":   plan.PlanName,
		"customer_id": customerID,
		"plan_items":  outplanitems,
	})
}
