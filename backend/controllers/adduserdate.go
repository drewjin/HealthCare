package controllers

import (
	"fmt"
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 为用户添加体检数据
// 仅允许机构用户添加
func AddUserDateByPlanID(ctx *gin.Context) {
	customerID := ctx.Param("customer_id")
	planID := ctx.Param("plan_id")

	type inputitem struct {
		ItemID    int     `json:"item_id" binding:"required"`  // 这里的item_id可以是PlanHeathItem的ID或HealthItem的ID
		ItemValue *string `json:"item_value" binding:"required"`
	}

	var input []inputitem

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// 验证用户权限 - 只有机构用户才能添加体检数据
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// 检查是否为机构用户 (UserType = 3)
	if user.UserType != 3 && user.UserType != 2 { // 不是机构用户也不是管理员
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Permission denied: only institution users can add health data",
		})
		return
	}

	//user_id,plan_id转换uint类型
	customerIDUint, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
		})
		return
	}
	planIDUint, err := strconv.ParseUint(planID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid plan ID",
		})
		return
	}

	// 验证套餐存在
	var plan models.Plan
	if err := global.DB.First(&plan, planIDUint).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Plan not found",
		})
		return
	}

	// 验证用户存在
	var customer models.User
	if err := global.DB.First(&customer, customerIDUint).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	// 验证用户是否已订阅该套餐
	var userPackage models.UserPackage
	if err := global.DB.Where("user_id = ? AND plan_id = ?", customerIDUint, planIDUint).First(&userPackage).Error; err != nil {
		// 如果用户没有订阅该套餐，自动创建订阅关系
		userPackage = models.UserPackage{
			UserID:        uint(customerIDUint),
			PlanID:        uint(planIDUint),
			InstitutionID: plan.RelationInstitutionID,
			Status:        0, // 0表示待检测
		}
		if err := global.DB.Create(&userPackage).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user package association",
			})
			return
		}
	}

	// 处理每个体检项目的数据
	for _, item := range input {
		if item.ItemValue == nil || *item.ItemValue == "" {
			continue // 跳过空值
		}

		// 获取健康项目ID（如果提供的是PlanHeathItem的ID，需要先查询对应的HealthItem ID）
		var healthItemID uint
		var planHeathItem models.PlanHeathItem
		
		// 尝试作为PlanHeathItem的ID查询
		if err := global.DB.First(&planHeathItem, item.ItemID).Error; err == nil {
			// 找到了匹配的PlanHeathItem
			healthItemID = planHeathItem.RelationHealthItemId
		} else {
			// 假设直接提供的是HealthItem的ID
			healthItemID = uint(item.ItemID)
		}

		// 检查此健康项目是否属于该套餐
		var count int64
		if err := global.DB.Model(&models.PlanHeathItem{}).
			Where("plan_id = ? AND health_item_id = ?", planIDUint, healthItemID).
			Count(&count).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to validate health item in plan",
			})
			return
		}

		if count == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Health item ID %d is not associated with plan ID %s", healthItemID, planID),
			})
			return
		}

		// 通过customerID、planID和healthItemID查找userhealthitem表中对应的记录
		var userHealthItem models.UserHealthItem
		result := global.DB.Where("user_id = ? AND plan_id = ? AND health_item_id = ?", 
			customerIDUint, planIDUint, healthItemID).First(&userHealthItem)
		
		// 如果不存在记录，则创建新记录
		if result.Error != nil {
			userHealthItem = models.UserHealthItem{
				RelationUserId:       uint(customerIDUint),
				RelationPlanId:       uint(planIDUint),
				RelationHealthItemId: healthItemID,
				ItemValue:            *item.ItemValue,
			}
			if err := global.DB.Create(&userHealthItem).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to create user health item: " + err.Error(),
				})
				return
			}
		} else {
			// 更新现有记录
			userHealthItem.ItemValue = *item.ItemValue
			if err := global.DB.Save(&userHealthItem).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to update user health item: " + err.Error(),
				})
				return
			}
		}
	}

	// 更新用户套餐状态为已完成
	if err := global.DB.Model(&models.UserPackage{}).
		Where("user_id = ? AND plan_id = ?", customerIDUint, planIDUint).
		Update("status", 1).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user package status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User health items updated successfully",
	})
}
