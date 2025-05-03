package controllers

import (
	"errors"
	"healthcare/controllers/utils"
	"healthcare/global"
	"healthcare/models"
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
