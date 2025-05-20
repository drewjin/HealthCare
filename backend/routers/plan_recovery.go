package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupPlanRecoveryRouter 设置套餐恢复相关路由
func SetupPlanRecoveryRouter(r *gin.Engine) {
	// 设置套餐恢复路由组
	planRecovery := r.Group("/api/admin/plans")
	planRecovery.Use(middlewares.AuthMiddleWare())
	planRecovery.Use(middlewares.AdminRequiredMiddleware())
	{
		// 获取已删除的套餐列表
		planRecovery.GET("/deleted", controllers.GetDeletedPlans)
		
		// 恢复已删除的套餐
		planRecovery.POST("/recover/:id", controllers.RecoverDeletedPlan)
	}
}
