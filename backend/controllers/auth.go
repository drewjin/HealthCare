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

func Register(ctx *gin.Context) {
	var user models.User

	// Bind the incoming JSON data to the user struct
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate user type
	if user.UserType < 1 || user.UserType > 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user type. Must be 1 (normal), 2 (admin), or 3 (institution)",
		})
		return
	}

	// Hash the password
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Password = hashedPwd

	// Generate a JWT token with user type
	token, err := utils.GenerateJWT(user.Username, user.UserType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create the user
	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send the JWT token, user ID, and user type
	ctx.JSON(http.StatusOK, gin.H{
		"token":     token,
		"uid":       user.ID,
		"user_type": user.UserType,
	})
}

func Login(ctx *gin.Context) {
	// Construct a struct to hold the input data
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User

	if err := global.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
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

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credentials",
		})
		return
	}

	// Generate token with user type
	token, err := utils.GenerateJWT(user.Username, user.UserType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := global.DB.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":     token,
		"uid":       user.ID,
		"user_type": user.UserType,
	})
}
