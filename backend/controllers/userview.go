package controllers

import (
	"encoding/json"
	"net/http"

	"healthcare/global"
	"healthcare/models"

	"github.com/gin-gonic/gin"
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
		InstitutionName string `json:"institution_name"`
		PlanName        string `json:"plan_name"`
		Items           string `json:"items"`
	}
	var records []record

	// 对每个套餐聚合健康记录
	for _, pkg := range userPackages {
		// 查询该用户该套餐下的所有健康项目
		var healthItems []models.UserHealthItem
		if err := global.DB.Where("user_id = ? AND plan_id = ?", userID, pkg.PlanID).Find(&healthItems).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve health items"})
			return
		}

		// 构建列表
		var list []map[string]string
		for _, hi := range healthItems {
			var itemName string
			_ = global.DB.Model(&models.HealthItem{}).Where("id = ?", hi.RelationHealthItemId).Select("item_name").First(&itemName).Error
			list = append(list, map[string]string{"item_name": itemName, "item_value": hi.ItemValue})
		}

		bytes, err := json.Marshal(list)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal items"})
			return
		}

		records = append(records, record{
			PlanID:          pkg.PlanID,
			InstitutionName: pkg.Institution.InstitutionName,
			PlanName:        pkg.Plan.PlanName,
			Items:           string(bytes),
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

	type output struct {
		PlanName  string `json:"plan_name"`
		ItemName  string `json:"item_name"`
		ItemValue string `json:"item_value"`
	}
	var outplanitems []output

	// 查询并返回
	var planItems []models.PlanHeathItem
	if err := global.DB.Where("plan_id = ?", input.PlanID).Find(&planItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan items"})
		return
	}

	for _, item := range planItems {
		var itemValue string
		_ = global.DB.Model(&models.UserHealthItem{}).Where("plan_id = ? AND health_item_id = ?", input.PlanID, item.RelationHealthItemId).Select("item_value").First(&itemValue).Error
		outplanitems = append(outplanitems, output{
			PlanName:  item.ThisPlan.PlanName,
			ItemName:  item.ThisHeathItem.ItemName,
			ItemValue: itemValue,
		})
	}

	c.JSON(http.StatusOK, gin.H{"plan_items": outplanitems})
}
