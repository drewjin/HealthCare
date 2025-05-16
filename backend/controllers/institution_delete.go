package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteInstitution 删除机构信息
func DeleteInstitution(ctx *gin.Context) {
	institutionID := ctx.Param("id")

	// 获取用户身份
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 查询机构是否存在
	var institution models.Institution
	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	// 验证权限：仅允许机构所有者或管理员删除机构
	if user.ID != institution.UserID && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to delete this institution",
		})
		return
	}

	// 开始事务
	tx := global.DB.Begin()
	
	// 先删除机构关联的套餐项目
	var plans []models.Plan
	if err := tx.Where("institution_id = ?", institution.ID).Find(&plans).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query institution plans: " + err.Error(),
		})
		return
	}
	
	// 获取套餐ID列表
	planIDs := make([]uint, 0, len(plans))
	for _, plan := range plans {
		planIDs = append(planIDs, plan.ID)
	}
	
	// 如果有套餐，删除套餐对应的检查项目关联
	if len(planIDs) > 0 {
		if err := tx.Where("plan_id IN ?", planIDs).Delete(&models.PlanHeathItem{}).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete plan health items: " + err.Error(),
			})
			return
		}
		
		// 删除套餐
		if err := tx.Where("institution_id = ?", institution.ID).Delete(&models.Plan{}).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete institution plans: " + err.Error(),
			})
			return
		}
	}
	
	// 最后删除机构
	if err := tx.Delete(&institution).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete institution: " + err.Error(),
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "机构及相关数据已成功删除",
	})
}
