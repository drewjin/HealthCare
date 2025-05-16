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
		institution.GET("", controllers.GetInstitutions)
		// Get User institutions
		institution.GET("/:id", middlewares.RequireUserType(3, 2), controllers.GetInstitutionDetail)
		// Create institution plans&items
		institution.POST("/:id/plans", middlewares.RequireUserType(3, 2), controllers.CreateInstitutionPlans)
		// Get institution Plans
		institution.GET("/:id/plans", controllers.GetInstitutionPlans)
		// Get institution Plan Items
		institution.POST("/:id/:plan_id/item", middlewares.RequireUserType(3, 2), controllers.CreateInstitutionPlans)
		// 更新机构相关信息
		institution.PATCH("/:id/update", middlewares.RequireUserType(3, 2), controllers.UpdateInsistutionPlanorItem)
		// 删除机构相关信息,删除套餐内一个体检项目

		// 删除套餐

		// 删除机构

	}
}
