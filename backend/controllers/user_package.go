package controllers

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SelectPackage allows users to select a package
func SelectPackage(ctx *gin.Context) {
	// Get user info from context
	username := ctx.GetString("username")
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var input struct {
		InstitutionID uint `json:"institution_id" binding:"required"`
		PlanID        uint `json:"plan_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Verify the institution and plan exist and the institution is approved
	var institution models.Institution
	if err := global.DB.First(&institution, input.InstitutionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Institution not found",
		})
		return
	}

	if institution.Status != 1 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "The institution is not approved",
		})
		return
	}

	var plan models.Plan
	if err := global.DB.Where("id = ? AND institution_id = ?", input.PlanID, input.InstitutionID).First(&plan).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Plan not found for this institution",
		})
		return
	}

	// Check if the user has already selected this package
	var existingSelection models.UserPackage
	result := global.DB.Where("user_id = ? AND institution_id = ? AND plan_id = ?",
		user.ID, input.InstitutionID, input.PlanID).First(&existingSelection)

	if result.Error == nil {
		// Already exists
		ctx.JSON(http.StatusConflict, gin.H{
			"message":      "You have already selected this package",
			"user_package": existingSelection,
		})
		return
	}

	// Create new user package selection
	newUserPackage := models.UserPackage{
		UserID:        user.ID,
		InstitutionID: input.InstitutionID,
		PlanID:        input.PlanID,
		Status:        0, // Default to pending
	}

	if err := global.DB.Create(&newUserPackage).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to select package: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Package selected successfully",
		"user_package": newUserPackage,
	})
}

// GetUserPackages returns all packages selected by a user
func GetUserPackages(ctx *gin.Context) {
	userID := ctx.Param("id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Get current user info
	username := ctx.GetString("username")
	var currentUser models.User
	if err := global.DB.Where("username = ?", username).First(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Only allow users to view their own packages unless they're an admin
	if currentUser.ID != uint(uid) && currentUser.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to view other user's packages",
		})
		return
	}

	var userPackages []models.UserPackage
	if err := global.DB.Where("user_id = ?", uid).
		Preload("Institution").
		Preload("Plan").
		Find(&userPackages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user packages: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_packages": userPackages,
	})
}
