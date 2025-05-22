package controllers

import (
	"HealthCare/backend/controllers/utils"
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

// 新建个人健康指标 oy
func CreateHealthItem(ctx *gin.Context) {
	var input struct {
		UserID         uint   `json:"user_id" binding:"required"`
		UserHealthInfo string `json:"user_health_info" binding:"required"`
	}

	// 绑定输入参数
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}

	fmt.Println(input)

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
	itemNameAndUserInfo := strings.Split(input.UserHealthInfo, "|")

	itemName := itemNameAndUserInfo[0]
	userInfo := itemNameAndUserInfo[1]

	// 执行插入操作
	result := tx.Exec(`INSERT INTO health_items (user_id, user_health_info, item_name) VALUES (?, ?, ?)`,
		input.UserID, userInfo, itemName)

	if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "新建个人健康档案指标失败: " + result.Error.Error(),
		})
		return
	}

	// 检查是否成功插入
	if result.RowsAffected == 0 {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "未能创建健康档案指标",
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
		"message": "新建个人健康档案指标成功",
		"data": gin.H{
			"user_id":          input.UserID,
			"user_health_info": input.UserHealthInfo,
		},
	})
}

// 删除个人健康指标 oy
func DelHealthItem(ctx *gin.Context) {
	ID := ctx.Param("id")

	// 参数验证
	if ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "用户ID不能为空",
		})
		return
	}
	id, err := utils.UnmarshalUint(ID)
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
	result := tx.Exec(`DELETE FROM health_items WHERE id = ?`, id)

	if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除健康档案指标失败: " + result.Error.Error(),
		})
		return
	}

	// 检查是否实际删除了记录
	if result.RowsAffected == 0 {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "未找到对应的健康档案记录",
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
		"message": "删除个人健康档案指标成功",
		"data": gin.H{
			"user_id": id,
		},
	})
}

// 修改个人健康指标 oy
func UpdateUserHealthItem(ctx *gin.Context) {
	var input struct {
		ID             uint   `json:"id" binding:"required"`
		UserHealthInfo string `json:"user_health_info" binding:"required"`
	}

	// 绑定输入参数
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// fmt.Print(input.ID, "\n")

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

	// fmt.Println("检查ID是否存在")
	// var originalInfo string
	// if err := tx.Raw("SELECT user_health_info FROM health_items WHERE id = ?", input.ID).Scan(&originalInfo).Error; err != nil {
	// 	tx.Rollback()
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询原始数据失败: " + err.Error()})
	// 	return
	// }
	// fmt.Printf("原始值为: %s\n", originalInfo)

	// // 执行更新操作
	// fmt.Println("修改")
	
	result := tx.Exec(`UPDATE health_items SET user_health_info = ? WHERE id = ?`,
		input.UserHealthInfo, input.ID)

	if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "修改健康档案指标失败: " + result.Error.Error(),
		})
		return
	}

	// 检查是否实际更新了记录
	if result.RowsAffected == 0 {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "未找到对应的健康档案记录",
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
		"message": "修改个人健康档案指标成功",
		"data": gin.H{
			"id":               input.ID,
			"user_health_info": input.UserHealthInfo,
		},
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
	var input struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Birthday string `json:"birthday"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Address  string `json:"address"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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

	result := tx.Exec(`UPDATE users SET username = ?,name=?,gender=?,birthday=?,phone=?,email=?,address=?  WHERE id = ?`,
		input.Username, input.Name, input.Gender, input.Birthday, input.Phone, input.Email, input.Address, userID)

	if result.Error != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "修改失败: " + result.Error.Error(),
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
