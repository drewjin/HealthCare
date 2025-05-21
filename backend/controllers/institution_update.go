package controllers

import (
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateInstitution 更新机构信息
func UpdateInstitution(ctx *gin.Context) {
	institutionID := ctx.Param("id")

	var input struct {
		InstitutionName          *string `json:"institution_name"`
		InstitutionPhone         *string `json:"institution_phone"`
		InstitutionAddress       *string `json:"institution_address"`
		InstitutionQualification *string `json:"institution_qualification"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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

	// 获取用户身份
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 验证权限：仅允许机构所有者或管理员更新机构信息
	if user.ID != institution.UserID && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to update this institution",
		})
		return
	}

	// 机构名称更新
	if input.InstitutionName != nil {
		// 检查机构名称唯一性，排除当前机构
		var count int64
		if err := global.DB.Model(&models.Institution{}).
			Where("institution_name = ? AND id != ?", *input.InstitutionName, institution.ID).
			Count(&count).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "数据库错误: " + err.Error(),
			})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "机构名称已存在",
			})
			return
		}

		institution.InstitutionName = *input.InstitutionName
	}

	// 更新其他字段
	if input.InstitutionPhone != nil && *input.InstitutionPhone != "" {
		institution.InstitutionPhone = *input.InstitutionPhone
	}

	if input.InstitutionAddress != nil && *input.InstitutionAddress != "" {
		institution.InstitutionAddress = *input.InstitutionAddress
	}

	if input.InstitutionQualification != nil && *input.InstitutionQualification != "" {
		institution.InstitutionQualification = *input.InstitutionQualification
	}

	// 保存更新
	if err := global.DB.Save(&institution).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新机构信息失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "机构信息更新成功",
		"institution": institution,
	})
}
