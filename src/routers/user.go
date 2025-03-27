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
		user.PUT("/users/:id/reset_pwd", controllers.ResetPwd)
		user.POST("/users/:id/relate")
	}
}
