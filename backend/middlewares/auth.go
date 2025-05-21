package middlewares

import (
	"HealthCare/backend/controllers/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 检查用户是否登录，是否有token，权限
func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Authorization Header",
			})
			ctx.Abort()
			return
		}
		username, userType, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token",
			})
			ctx.Abort()
			return
		}
		ctx.Set("username", username)
		ctx.Set("user_type", userType)
		ctx.Next()
	}
}
