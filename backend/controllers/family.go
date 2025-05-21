package controllers

import (
	"HealthCare/backend/controllers/utils"
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatFamily(ctx *gin.Context) {
	var input struct {
		RelativeUsername string `json:"relative_username"`
		Relationship     string `json:"relationship"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := ctx.Param("id")
	thisUserID, err := utils.UnmarshalUint(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var relative models.User
	if err := global.DB.Where("username = ?", input.RelativeUsername).First(&relative).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Credentials",
			})
		}
		return
	}

	if err := utils.CreateFamilyRequest(thisUserID, relative.ID, relative.Email, input.Relationship); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

// GetPendingFamilyRequests 获取待处理的家庭关系请求
func GetPendingFamilyRequests(ctx *gin.Context) {
	userID := ctx.Param("id")
	thisUserID, err := utils.UnmarshalUint(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var requests []models.Family
	if err := global.DB.Preload("ThisUser").Where("relative_id = ? AND status = 0", thisUserID).Find(&requests).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending requests"})
		return
	}

	var response []gin.H
	for _, req := range requests {
		response = append(response, gin.H{
			"id":           req.ID,
			"requester":    req.ThisUser.Username,
			"name":         req.ThisUser.Name,
			"relationship": req.Relationship,
			"created_at":   req.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// HandleFamilyRequest 处理家庭关系请求
func HandleFamilyRequest(ctx *gin.Context) {
	var input struct {
		Accept bool `json:"accept"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID := ctx.Param("requestId")
	reqID, err := utils.UnmarshalUint(requestID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.Param("id")
	thisUserID, err := utils.UnmarshalUint(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request models.Family
	if err := global.DB.Where("id = ? AND relative_id = ? AND status = 0", reqID, thisUserID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch request"})
		}
		return
	}

	if input.Accept {
		request.Status = 1
		if err := global.DB.Save(&request).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Request accepted successfully"})
	} else {
		if err := global.DB.Delete(&request).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Request rejected successfully"})
	}
}

// GetConfirmedFamilyMembers 获取已确认的家庭关系
func GetConfirmedFamilyMembers(ctx *gin.Context) {
	userID := ctx.Param("id")
	thisUserID, err := utils.UnmarshalUint(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var relationships []models.Family
	if err := global.DB.Preload("ThisUser").Preload("Relative").
		Where("(user_id = ? OR relative_id = ?) AND status = 1", thisUserID, thisUserID).
		Find(&relationships).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch family relationships"})
		return
	}

	var response []gin.H
	for _, rel := range relationships {
		var member gin.H
		var relationship string
		if rel.UserID == thisUserID {
			member = gin.H{
				"username": rel.Relative.Username,
				"user_id":  rel.ThisUser.ID,
				"id":       rel.Model.ID,
				"rid":      rel.ID,
				"name":     rel.Relative.Name,
			}
			relationship = rel.Relationship
		} else {
			member = gin.H{
				"username": rel.ThisUser.Username,
				"user_id":  rel.ThisUser.ID,
				"id":       rel.Model.ID,
				"rid":      rel.ID,
				"name":     rel.ThisUser.Name,
			}
			relationship = reverseRelationship(rel.Relationship)
		}
		response = append(response, gin.H{
			"username":     member["username"],
			"name":         member["name"],
			"user_id":      member["user_id"],
			"id":           member["id"],
			"rid":          member["rid"],
			"relationship": relationship,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// reverseRelationship 反转关系类型
func reverseRelationship(relationship string) string {
	switch relationship {
	case "父亲":
		return "儿子"
	case "母亲":
		return "儿子"
	case "儿子":
		return "父亲" // 这里可能需要根据性别判断是父亲还是母亲
	case "女儿":
		return "父亲" // 这里可能需要根据性别判断是父亲还是母亲
	case "配偶":
		return "配偶"
	default:
		return relationship
	}
}

// 删除授权状态 oy
func DelFamilyStatus(ctx *gin.Context) {
	rid := ctx.Param("id")
	id, err := utils.UnmarshalUint(rid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "服务器内部错误",
			})
		}
	}()

	// 执行删除操作
	result := tx.Exec(`DELETE FROM families WHERE id = ?`, id)
	if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除授权状态失败: " + result.Error.Error(),
		})
		return
	}

	// 检查是否实际删除了记录
	if result.RowsAffected == 0 {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "未找到对应的授权记录",
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "提交事务失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "删除授权状态成功",
	})
}

// 修改亲友关系名字 oy
func UpdateFamilyName(ctx *gin.Context) {
	var input struct {
		ID           uint   `json:"id" binding:"required"`
		Relationship string `json:"relationship" binding:"required"`
	}

	// 绑定输入参数
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "服务器内部错误",
			})
		}
	}()

	// 执行更新操作
	result := tx.Exec(`UPDATE families SET relationship = ? WHERE id = ?`,
		input.Relationship, input.ID)

	if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "修改亲友关系失败: " + result.Error.Error(),
		})
		return
	}

	// 检查是否实际更新了记录
	if result.RowsAffected == 0 {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "未找到对应的亲友记录",
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "提交事务失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "修改亲友关系成功",
	})
}
