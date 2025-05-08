package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(routers *gin.Engine) {
	user := routers.Group("/api")
	user.Use(middlewares.AuthMiddleWare())
	{
		user.GET("/users/:id/profile", controllers.GetUserProfileByID)
		user.PUT("/users/:id/reset_pwd", controllers.ResetPwd)
	}
}
