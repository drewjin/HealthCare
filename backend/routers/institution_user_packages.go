package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupInstitutionUserPackagesRouter 设置机构用户套餐相关路由
func SetupInstitutionUserPackagesRouter(r *gin.Engine) {
	// 获取机构下的用户套餐列表
	institutionUserPackages := r.Group("/api/institution/user-packages")
	institutionUserPackages.Use(middlewares.AuthMiddleWare())
	{
		institutionUserPackages.GET("", controllers.GetUserPackagesByInstitution)
	}
}
