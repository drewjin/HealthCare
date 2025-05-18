package controllers

import (
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
		ItemID    int     `json:"item_id" binding:"required"`
		ItemValue *string `json:"item_value" binding:"required"`
	}

	var input []inputitem

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// 通过customerID和planID查找userhealthitem表中对应的记录
	var userHealthItems []models.UserHealthItem
	if err := global.DB.Where("user_id = ? AND plan_id = ?", customerID, planID).Find(&userHealthItems).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user health items",
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

	// 如果没有记录，则插入新的记录
	if len(userHealthItems) == 0 {
		for _, item := range input {
			newItem := models.UserHealthItem{
				RelationUserId:       uint(customerIDUint),
				RelationPlanId:       uint(planIDUint),
				RelationHealthItemId: uint(item.ItemID),
				ItemValue:            *item.ItemValue,
			}
			if err := global.DB.Create(&newItem).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to create user health item",
				})
				return
			}
		}
	} else {
		// 更新userhealthitem表中的记录
		for _, item := range input {
			if item.ItemValue != nil && *item.ItemValue != "" {
				if err := global.DB.Model(&models.UserHealthItem{}).
					Where("user_id = ? AND plan_id = ? AND health_item_id = ?", customerID, planID, item.ItemID).
					Update("item_value", item.ItemValue).Error; err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to update user health items",
					})
					return
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User health items updated successfully",
	})

}
