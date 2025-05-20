package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupHealthItemsRouter 设置健康检查项目相关路由
func SetupHealthItemsRouter(r *gin.Engine) {
	healthItems := r.Group("api/healthitems")
	healthItems.Use(middlewares.AuthMiddleWare())
	{
		// 获取所有健康检查项目
		healthItems.GET("", controllers.GetAllHealthItems)
		
		// 获取指定ID的健康检查项目
		healthItems.GET("/:id", controllers.GetHealthItemByID)
		
		// 更新健康检查项目
		healthItems.PATCH("/:id", controllers.UpdateHealthItem)
		
		// 更新套餐中项目的描述
		healthItems.PATCH("/plan-item", controllers.UpdatePlanItemDescription)
	}
}
