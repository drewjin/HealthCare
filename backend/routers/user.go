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
		// Add new endpoint to get institution by user ID
		user.GET("/user/:id/institution", controllers.GetInstitutionByUserId)
	}
}
