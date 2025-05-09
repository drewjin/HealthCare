package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateInstitution(ctx *gin.Context) {
	var institution models.Institution
	if err := ctx.ShouldBindJSON(&institution); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get user ID from the path
	userID := ctx.Param("id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Verify that the user is an institution user
	var user models.User
	if err := global.DB.First(&user, uid).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.UserType != 3 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only institution users can create institution information",
		})
		return
	}

	// Check if institution already exists for this user
	var existingInstitution models.Institution
	err = global.DB.Where("user_id = ?", uid).First(&existingInstitution).Error
	if err == nil {
		// 找到了已存在的机构记录
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Institution already exists for this user",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 发生了除了"记录未找到"之外的错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 如果是 ErrRecordNotFound 错误，说明可以创建新机构，继续执行

	institution.UserID = uint(uid)
	institution.Status = 0 // Set status as pending

	fmt.Printf("Creating institution: %+v\n", institution)

	if err := global.DB.Create(&institution).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Institution information submitted successfully",
	})
}

func GetPendingInstitutions(ctx *gin.Context) {
	// Verify that the user is an admin
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin users can view pending institutions",
		})
		return
	}

	var institutions []models.Institution
	if err := global.DB.Where("status = ?", 0).Find(&institutions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, institutions)
}

func ReviewInstitution(ctx *gin.Context) {
	// Verify that the user is an admin
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin users can review institutions",
		})
		return
	}

	institutionID := ctx.Param("id")
	var input struct {
		Approved bool `json:"approved"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var institution models.Institution
	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	institution.Status = map[bool]uint8{true: 1, false: 2}[input.Approved]
	if err := global.DB.Save(&institution).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": map[bool]string{
			true:  "Institution approved successfully",
			false: "Institution rejected",
		}[input.Approved],
	})
}

func GetAllApprovedInstitutions(ctx *gin.Context) {
	var institutions []models.Institution
	if err := global.DB.Where("status = ?", 1).Find(&institutions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, institutions)
}

func GetInstitutionDetail(ctx *gin.Context) {
	institutionID := ctx.Param("id")
	var institution models.Institution

	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	// 检查请求是否来自管理员
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 如果不是管理员，只能看到已批准的机构
	if user.UserType != 2 && institution.Status != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to view this institution",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"institution": institution,
		"isAdmin":     user.UserType == 2,
	})
}

func GetInstitutionPackages(ctx *gin.Context) {
	institutionID := ctx.Param("id")
	var institution models.Institution

	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	// 只允许查看已批准的机构的套餐
	if institution.Status != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "This institution is not approved",
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

	// 假设ExaminationPackage字段存储了JSON格式的套餐信息
	// 实际应用中，这里可能需要解析JSON或从其他表中查询

	ctx.JSON(http.StatusOK, gin.H{
		"institution": institution,
		"isAdmin":     user.UserType == 2,
		"packages":    institution.ExaminationPackage,
	})
}

func UpdateInstitutionPackages(ctx *gin.Context) {
	institutionID := ctx.Param("id")
	var institution models.Institution

	if err := global.DB.First(&institution, institutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	// 验证操作者权限（仅机构所有者或管理员可以更新套餐）
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// 检查是否为机构所有者或管理员
	if user.ID != institution.UserID && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to update packages for this institution",
		})
		return
	}

	// 解析请求体
	var input struct {
		Packages string `json:"packages"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 验证JSON格式
	var packageTest interface{}
	if err := json.Unmarshal([]byte(input.Packages), &packageTest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format for packages",
		})
		return
	}

	// 更新套餐信息
	institution.ExaminationPackage = input.Packages
	if err := global.DB.Save(&institution).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Examination packages updated successfully",
	})
}
