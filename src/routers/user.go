package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(router *gin.Engine) {
	user := router.Group("/api")
	user.Use(middlewares.AuthMiddleWare())
	{
		user.GET("/users/:id/profile", controllers.GetUserProfileByID)
		user.PUT("/reset-pwd", controllers.ResetPwd)
	}
}
