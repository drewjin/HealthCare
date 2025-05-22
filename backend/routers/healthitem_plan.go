package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupHealthItemPlanRouter 设置健康项目与套餐关联的路由
func SetupHealthItemPlanRouter(r *gin.Engine) {
	healthItemPlan := r.Group("api/healthitem-plan")
	healthItemPlan.Use(middlewares.AuthMiddleWare())
	{
		// 获取套餐关联的健康项目
		healthItemPlan.GET("/plans/:plan_id/items", controllers.GetPlanHealthItems)

		// 将健康项目关联到套餐
		healthItemPlan.POST("/plans/:plan_id/items", middlewares.RequireUserType(3, 2), controllers.AssociatePlanHealthItems)

		// 添加单个健康项目到套餐
		healthItemPlan.POST("/plans/:plan_id/item", middlewares.RequireUserType(3, 2), controllers.AddHealthItemToPlan)

		// 获取用户特定套餐的健康数据
		healthItemPlan.GET("/users/:user_id/plans/:plan_id", controllers.GetUserPlanHealthItems)

		// 按机构获取用户的健康数据
		healthItemPlan.GET("/users/:user_id/institutions/:institution_id", controllers.GetUserInstitutionHealthItems)
	}
}
