package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 查看该用户下的所有体检项目(包含所有体检套餐)
func GetAllItems(ctx *gin.Context) {
	username := ctx.GetString("username")
	// 通过username查找user表中对应的userid
	var userID uint
	if err := global.DB.Model(&models.User{}).Where("username = ?", username).Select("id").First(&userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	type output struct {
		ItemName  string `json:"item_name"`
		ItemValue string `json:"item_value"`
	}

	// 查找userhealthitem表中username对应userid的记录
	var userHealthItems []models.UserHealthItem
	if err := global.DB.Where("user_id = ?", userID).Find(&userHealthItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user health items"})
		return
	}

	// 去除重复的体检项目，只保留最新记录
	uniqueItems := make(map[uint]models.UserHealthItem)
	for _, item := range userHealthItems {
		if existingItem, exists := uniqueItems[item.RelationHealthItemId]; !exists || item.UpdatedAt.After(existingItem.UpdatedAt) {
			uniqueItems[item.RelationHealthItemId] = item
		}
	}

	// 将体检项目转换成output结构体
	var items []output
	for _, item := range uniqueItems {
		var itemName string
		if err := global.DB.Model(&models.HealthItem{}).Where("id = ?", item.RelationHealthItemId).Select("item_name").First(&itemName).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item name"})
			return
		}
		items = append(items, output{
			ItemName:  itemName,
			ItemValue: item.ItemValue,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"items": items})
}

// 查看指定套餐下的体检项目
func GetItemsByPlanID(c *gin.Context) {
	var input struct {
		PlanID uint `json:"plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 通过plan_id查找planheathitem中对应的体检项目
	var planItems []models.PlanHeathItem
	if err := global.DB.Where("plan_id = ?", input.PlanID).Find(&planItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan items"})
		return
	}

	// 根据planid，itemid查找userhealthitem表中对应的记录，并转换成output结构体
	var outplanitems []interface{}

	for _, item := range planItems {
		var itemName string
		if err := global.DB.Model(&models.HealthItem{}).Where("id = ?", item.RelationHealthItemId).Select("item_name").First(&itemName).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item name"})
			return
		}

		var itemValue string
		if err := global.DB.Model(&models.UserHealthItem{}).Where("plan_id = ? AND health_item_id = ?", input.PlanID, item.RelationHealthItemId).Select("item_value").First(&itemValue).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item value"})
			return
		}

		outplanitems = append(outplanitems, struct {
			PlanName  string `json:"plan_name"`
			ItemName  string `json:"item_name"`
			ItemValue string `json:"item_value"`
		}{
			PlanName:  item.ThisPlan.PlanName,
			ItemName:  itemName,
			ItemValue: itemValue,
		})
	}
}
