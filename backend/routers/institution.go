package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupInstitutionRouter(r *gin.Engine) {
	institution := r.Group("/api/institutions")
	institution.Use(middlewares.AuthMiddleWare())
	{
		// Only institution users can create institution info
		institution.POST("/:id", middlewares.RequireUserType(3), controllers.CreateInstitution)
		// Only admin users can view pending and review institutions
		institution.GET("/pending", middlewares.RequireUserType(2), controllers.GetPendingInstitutions)
		institution.POST("/:id/review", middlewares.RequireUserType(2), controllers.ReviewInstitution)

		// New APIs for institution viewing
		institution.GET("", controllers.GetAllApprovedInstitutions)          // Get all approved institutions
		institution.GET("/:id", controllers.GetInstitutionDetail)            // Get institution details
		institution.GET("/:id/packages", controllers.GetInstitutionPackages) // Get institution examination packages

		// API for updating institution packages (only institution owner or admin)
		institution.PUT("/:id/packages", controllers.UpdateInstitutionPackages)
	}
}
