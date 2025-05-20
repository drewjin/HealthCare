package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupUserPackageRouter 设置用户套餐选择相关路由
func SetupUserPackageRouter(r *gin.Engine) {
	// 用户选择套餐
	userPackage := r.Group("/api/user-packages")
	userPackage.Use(middlewares.AuthMiddleWare())
	{
		userPackage.POST("", controllers.SelectPackage)
		userPackage.GET("/:id", controllers.GetUserPackages)
	}
}
