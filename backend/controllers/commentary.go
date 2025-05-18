package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 发布评论
func AddCommentary(ctx *gin.Context) {

	// 从请求中获取评论内容
	var input struct {
		RelationUserId uint   `json:"user_id"`
		RelationPlanId uint   `json:"plan_id"`
		Commentary     string `json:"commentary"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	var commentary models.Commentary
	commentary.RelationUserId = input.RelationUserId
	commentary.RelationPlanId = input.RelationPlanId
	commentary.Commentary = input.Commentary

	// 添加记录到commentary表
	if err := global.DB.Create(&commentary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add commentary",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Commentary added successfully",
	})
}

// 删除评论
func DeleteCommentary(ctx *gin.Context) {
	commentaryID := ctx.Param("id")
	var commentary models.Commentary
	if err := global.DB.Where("id = ?", commentaryID).First(&commentary).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Commentary not found",
		})
		return
	}

	username := ctx.GetString("username")
	// 检查用户是否有权限删除评论
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 只有评论的创建者或管理员可以删除评论
	if commentary.RelationUserId != user.ID && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to delete this commentary",
		})
		return
	}

	if err := global.DB.Unscoped().Delete(&commentary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete commentary",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Commentary deleted successfully",
	})
}

// 查看评论(套餐id)
func GetCommentaryByPlanID(ctx *gin.Context) {
	planID := ctx.Param("id")
	var commentaries []models.Commentary
	if err := global.DB.Where("plan_id = ?", planID).Find(&commentaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve commentaries",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"commentaries": commentaries,
	})
}

// 查看评论(用户id)
func GetCommentaryByUserID(ctx *gin.Context) {
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var commentaries []models.Commentary
	if err := global.DB.Where("user_id = ?", user.ID).Find(&commentaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve commentaries",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"commentaries": commentaries,
	})
}
