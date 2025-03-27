package routers

import (
	"healthcare/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		// TODO: Implement other authentication routes here
		// EMail and Phone verification
	}
}
