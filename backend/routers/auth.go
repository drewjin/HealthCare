package routers

import (
	"HealthCare/backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(routers *gin.Engine) {
	auth := routers.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		// TODO: Implement other authentication routes here
		// EMail and Phone verification
	}
}
