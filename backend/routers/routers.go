package routers

import (
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	routers := gin.Default()

	routers.Use(middlewares.SetupCorsMiddleware())
	{
		SetupAuthRouter(routers)
		SetupUserRouter(routers)
		SetupFamilyRouter(routers)
		SetupInstitutionRouter(routers)
		SetupCommentaryRouter(routers)
		SetupUserViewRouter(routers)
	}
	return routers
}
