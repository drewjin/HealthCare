package controllers

import (
	"HealthCare/backend/controllers/utils"
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllHealthItems 获取所有健康检查项目
func GetAllHealthItems(ctx *gin.Context) {
	var healthItems []models.HealthItem

	// 支持关键词搜索
	keyword := ctx.Query("keyword")
	if keyword != "" {
		if err := global.DB.Where("item_name LIKE ?", "%"+keyword+"%").Find(&healthItems).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取健康检查项目失败: " + err.Error(),
			})
			return
		}
	} else {
		if err := global.DB.Find(&healthItems).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取健康检查项目失败: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": healthItems,
	})
}

// oy
func GetAllHealthItemsList(ctx *gin.Context) {
	var items []map[string]interface{}

	global.DB.
		Table("health_items").
		Scan(&items)

	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
	})

}
func GetAllHealthItemsByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	thisUserID, _ := utils.UnmarshalUint(userID)

	var items []map[string]interface{}

	rawSQL := "SELECT * FROM health_items WHERE user_id = ?"
	global.DB.Raw(rawSQL, thisUserID).Scan(&items)
	// fmt.Printf("items: %v\n", items)
	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

// GetHealthItemByID 获取指定ID的健康检查项目
func GetHealthItemByID(ctx *gin.Context) {
	itemID := ctx.Param("id")

	var healthItem models.HealthItem
	if err := global.DB.First(&healthItem, itemID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "健康检查项目不存在",
		})
		return
	}

	// 获取关联的套餐信息
	var plans []struct {
		PlanID          uint    `json:"PlanID"`
		PlanName        string  `json:"PlanName"`
		ItemDescription string  `json:"ItemDescription"`
		PlanPrice       float64 `json:"PlanPrice"`
	}

	global.DB.Table("plan_heath_items").
		Select("plan_heath_items.plan_id as PlanID, plans.plan_name as PlanName, plan_heath_items.item_description as ItemDescription, plans.plan_price as PlanPrice").
		Joins("JOIN plans ON plans.id = plan_heath_items.plan_id").
		Where("plan_heath_items.health_item_id = ?", itemID).
		Find(&plans)

	ctx.JSON(http.StatusOK, gin.H{
		"item":  healthItem,
		"plans": plans,
	})
}

// UpdateHealthItem 更新健康检查项目
func UpdateHealthItem(ctx *gin.Context) {
	itemID := ctx.Param("id")

	var input struct {
		ItemName string `json:"item_name"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 更新健康检查项目
	if err := global.DB.Model(&models.HealthItem{}).Where("id = ?", itemID).Update("item_name", input.ItemName).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新健康检查项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "健康检查项目已更新",
	})
}

// UpdatePlanItemDescription 更新套餐中项目的描述
func UpdatePlanItemDescription(ctx *gin.Context) {
	var input struct {
		PlanID          uint   `json:"plan_id" binding:"required"`
		ItemID          uint   `json:"item_id" binding:"required"`
		ItemDescription string `json:"item_description"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// 更新套餐项目描述
	if err := global.DB.Model(&models.PlanHeathItem{}).
		Where("plan_id = ? AND health_item_id = ?", input.PlanID, input.ItemID).
		Update("item_description", input.ItemDescription).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新项目描述失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "项目描述已更新",
	})
}
