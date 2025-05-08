package middlewares

import (
	"healthcare/global"
	"healthcare/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireUserType(allowedTypes ...uint8) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")
		var user models.User
		if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			ctx.Abort()
			return
		}

		allowed := false
		for _, t := range allowedTypes {
			if user.UserType == t {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Unauthorized access: insufficient privileges",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
