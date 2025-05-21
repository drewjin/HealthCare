package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

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
		// 获取用户的机构列表
		institution.GET("", controllers.GetInstitutions)
		// 获取用户的机构列表
		institution.GET("/:id", middlewares.RequireUserType(3, 2, 1), controllers.GetInstitutionDetail)
		// 为指定机构添加套餐
		institution.POST("/:id/plans", middlewares.RequireUserType(3, 2), controllers.CreateInstitutionPlans)
		// 获取指定机构的套餐列表
		institution.GET("/:id/plans", controllers.GetInstitutionPlans)
		// 创建机构下的指定套餐的体检项目
		institution.POST("/:id/:plan_id/item", middlewares.RequireUserType(3, 2), controllers.CreateInstitutionPlans)
		// 更新机构相关信息
		institution.PATCH("/:id/update", middlewares.RequireUserType(3, 2), controllers.UpdateInstitution)
		// 更新套餐的体检项目信息
		institution.PATCH("/:id/item", middlewares.RequireUserType(3, 2), controllers.UpdateInstitutionPlanorItem)

		// 删除都是物理删除
		// 删除套餐或检查项目信息,删除套餐内一个体检项目
		institution.DELETE("/plan/item", middlewares.RequireUserType(3, 2), controllers.DeleteInsistutionPlanonItem)
		// 删除套餐
		institution.DELETE("/plan", middlewares.RequireUserType(3, 2), controllers.DeleteInstitutionPlan)
		// 删除机构
		institution.DELETE("/:id", middlewares.RequireUserType(3, 2), controllers.DeleteInstitution)
	}
}
