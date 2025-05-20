package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupHealthItemManagerRouter 设置健康项目管理相关路由
func SetupHealthItemManagerRouter(r *gin.Engine) {
	healthItemManager := r.Group("api/healthitem-manager")
	healthItemManager.Use(middlewares.AuthMiddleWare())
	{
		// 创建健康项目模板
		healthItemManager.POST("/template", controllers.CreateHealthItemTemplate)
		
		// 保存健康项目模板到数据库
		healthItemManager.POST("/save-template", controllers.SaveHealthItemTemplate)
		
		// 更新健康项目值
		healthItemManager.PUT("/values", controllers.UpdateHealthItemValues)
		
		// 获取健康项目解析后的键值对
		healthItemManager.GET("/values/:id", controllers.GetHealthItemValues)
	}
}