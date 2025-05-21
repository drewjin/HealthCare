package routers

import (
	"HealthCare/backend/middlewares"

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
		SetupAddUserDataRouter(routers)
		SetupImageOcrRouter(routers)
		SetupUserPlansRouter(routers)
		SetupHealthItemsRouter(routers)
		SetupHealthItemManagerRouter(routers)
		SetupHealthItemPlanRouter(routers)
		SetupInstitutionUserPackagesRouter(routers)
		SetupPlanRecoveryRouter(routers)
		SetupUserPackageStatusRouter(routers)
		SetupUserPackageRouter(routers)
	}
	return routers
}
