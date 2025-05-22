package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupFamilyRouter(r *gin.Engine) {
	family := r.Group("/api/family")
	family.Use(middlewares.AuthMiddleWare())
	{
		family.POST("/request/:id", controllers.CreatFamily)
		family.GET("/pending/:id", controllers.GetPendingFamilyRequests)
		family.POST("/handle/:id/:requestId", controllers.HandleFamilyRequest)
		family.GET("/confirmed/:id", controllers.GetConfirmedFamilyMembers) // 新增获取已确认家庭关系的路由
		family.GET("/del_confirmed/:id", controllers.DelFamilyStatus)       // oy
		family.POST("/update_family_name", controllers.UpdateFamilyName)    //oy

	}
}
