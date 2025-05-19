package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserPlansRouter(r *gin.Engine) {
	// 获取用户套餐列表
	userPlans := r.Group("/api/users/:id/plans")
	userPlans.Use(middlewares.AuthMiddleWare())
	{
		userPlans.GET("", controllers.GetUserPlans)
	}

	// 获取套餐项目列表
	planItems := r.Group("/api/plans")
	planItems.Use(middlewares.AuthMiddleWare())
	{
		planItems.GET("/:id/items", controllers.GetPlanItems)
		planItems.GET("/:id", controllers.GetPlanDetails)
	}
}
