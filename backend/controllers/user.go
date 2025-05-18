package controllers

import (
	"errors"
	"fmt"
	"healthcare/controllers/utils"
	"healthcare/global"
	"healthcare/models"
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

// 用户更新个人信息
func UpdateUserProfile(ctx *gin.Context) {
	userID := ctx.Param("id")
	var user models.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var input struct {
		Username *string `json:"username"`
		Name     *string `json:"name"`
		Gender   *string `json:"gender"`
		Birthday *string `json:"birthday"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		Address  *string `json:"address"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if input.Username != nil && *input.Username != "" {
		updates["username"] = *input.Username
	}
	if input.Name != nil && *input.Name != "" {
		updates["name"] = *input.Name
	}
	if input.Gender != nil && *input.Gender != "" {
		updates["gender"] = *input.Gender
	}
	if input.Birthday != nil && *input.Birthday != "" {
		updates["birthday"] = *input.Birthday
	}
	if input.Phone != nil && *input.Phone != "" {
		updates["phone"] = *input.Phone
	}
	if input.Email != nil && *input.Email != "" {
		updates["email"] = *input.Email
	}
	if input.Address != nil && *input.Address != "" {
		updates["address"] = *input.Address
	}

	fmt.Println(updates)

	switch len(updates) {
	case 0:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields to update",
		})
		return
	default:
		for key, value := range updates {
			result := global.DB.Model(&user).Update(key, value)
			if result.Error != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": result.Error.Error(),
					"field": key,
					"value": value,
				})
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
	})

}

// 删除(注销)用户
func DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var user models.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查是本人还是管理员
	username := ctx.GetString("username")
	if username != user.Username && user.UserType != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to delete this user",
		})
		return
	}

	if err := global.DB.Unscoped().Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"logout":  true,
	})
}

// 管理员更改用户权限
func UpdateUserPermission(ctx *gin.Context) {
	changeuserID := ctx.Param("id")

	var input struct {
		UserType uint8 `json:"user_type"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	utils.UpdateIt(&models.User{}, changeuserID, "user_type", input.UserType)

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "User permission updated successfully",
		"user_id":   changeuserID,
		"user_type": input.UserType,
		"logout":    true,
	})
}
