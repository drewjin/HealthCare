package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupUserPackageStatusRouter 设置用户套餐状态更新相关路由
func SetupUserPackageStatusRouter(r *gin.Engine) {
	// 更新用户套餐状态
	userPackageStatus := r.Group("/api/user-packages/:user_id/:plan_id/status")
	userPackageStatus.Use(middlewares.AuthMiddleWare())
	{
		userPackageStatus.PATCH("", controllers.UpdateUserPackageStatus)
	}
}
