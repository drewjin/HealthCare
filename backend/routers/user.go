package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(routers *gin.Engine) {
	user := routers.Group("/api/users")
	user.Use(middlewares.AuthMiddleWare())
	{
		// 查看用户个人信息
		user.GET("/:id/profile", controllers.GetUserProfileByID)
		// 新建个人健康指标 oy
		user.POST("/create_health_item", controllers.CreateHealthItem)
		// 删除个人健康指标 oy
		user.GET("/:id/del_health_item", controllers.DelHealthItem)
		// 修改个人健康指标 oy
		user.POST("/update_use_health_item", controllers.UpdateUserHealthItem)
		// 用户重设密码
		user.PUT("/:id/reset_pwd", controllers.ResetPwd)
		// 查找用户管理的机构
		user.GET("/:id/institution", controllers.GetInstitutionByUserId)
		// 用户更新个人信息
		user.POST("/:id/profile", controllers.UpdateUserProfile)
		// 删除用户(物理删除，前端二次确认不予后悔)
		user.DELETE("/:id", middlewares.RequireUserType(2, 1), controllers.DeleteUser)
		// 管理员更改用户权限
		user.PATCH("/:id/permission", middlewares.RequireUserType(2), controllers.UpdateUserPermission)

		// 用户选择套餐
		user.POST("/packages", controllers.SelectPackage)
		// 用户查看已经选择的套餐
		user.GET("/:id/packages", controllers.GetUserPackages)
	}
}
