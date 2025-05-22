package middlewares

import (
	"HealthCare/backend/global"
	"HealthCare/backend/models"
	"net/http"

	"slices"

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

		allowed := slices.Contains(allowedTypes, user.UserType)

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
