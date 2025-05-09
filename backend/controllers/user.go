package controllers

import (
	"errors"
	"healthcare/global"
	"healthcare/models"
	"healthcare/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserProfileByID(ctx *gin.Context) {
	user_id := ctx.Param("id")

	var user models.User
	if err := global.DB.Select("id, username, name, email, phone, gender, birthday, address, user_type").
		Where("id = ?", user_id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	profile := struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Gender   string `json:"gender"`
		Birthday string `json:"birthday"`
		Address  string `json:"address"`
		UserType uint8  `json:"user_type"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Gender:   user.Gender,
		Birthday: user.Birthday,
		Address:  user.Address,
		UserType: user.UserType,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

func ResetPwd(ctx *gin.Context) {
	// Construct a struct to hold the input data
	var input struct {
		PrevPassword       string `json:"prev_password"`
		NewPassword        string `json:"new_password"`
		NewPasswordConfirm string `json:"new_password_confirm"`
	}

	// Bind the JSON request body into the input struct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Getting the user ID from the context
	userID := ctx.Param("id")

	// Query the user record from the database
	var user models.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	// Check if the previous password matches the user's password
	if !utils.CheckPassword(input.PrevPassword, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credentials: password does not match",
		})
		return
	}

	// Check if the new passwords match
	if input.NewPassword != input.NewPasswordConfirm {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credentials: new passwords do not match",
		})
		return
	}

	// Hash the new password
	hashedPwd, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = hashedPwd

	// Saving the updated user record
	result := global.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPwd)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Check if the password was updated
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found or no changes made",
		})
		return
	}

	// Clean up the client token cookie
	ctx.SetCookie("token", "", -1, "/", "", false, true)

	// Return a success message
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. Please login again.",
		"logout":  true, // Frontend should redirect to login page
	})
}

func GetInstitutionByUserId(ctx *gin.Context) {
	userID := ctx.Param("id")

	var institution models.Institution
	if err := global.DB.Where("user_id = ?", userID).First(&institution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No institution found for this user",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, institution)
}
