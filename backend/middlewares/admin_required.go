package middlewares

import (
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminRequiredMiddleware 要求用户是管理员
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")

		var user models.User
		if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的用户",
			})
			c.Abort()
			return
		}

		// 检查是否为管理员 (UserType = 2)
		if user.UserType != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
