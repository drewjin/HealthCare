package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserViewRouter(r *gin.Engine) {
	userview := r.Group("api/userview")
	userview.Use(middlewares.AuthMiddleWare())
	{
		// 查看所有体检项目(包含所有体检套餐)
		userview.GET("/", controllers.GetAllItems)
		// 查看指定套餐体检项目
		userview.GET("/plan", controllers.GetItemsByPlanID)
	}
}
