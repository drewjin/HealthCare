package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCommentaryRouter(routers *gin.Engine) {
	commentary := routers.Group("/api/commentary")
	commentary.Use(middlewares.AuthMiddleWare())
	{
		// oy
		commentary.GET("/get/user_list", controllers.GetCommentaryList)
		// 发布评论
		commentary.POST("/add", controllers.AddCommentary)
		// 删除评论(物理删除)
		commentary.DELETE("/delete/:id", controllers.DeleteCommentary)
		// 查看评论(套餐id)
		commentary.GET("/get/plan/:id", controllers.GetCommentaryByPlanID)
		// 查看评论(用户id)
		commentary.GET("/get/user", controllers.GetCommentaryByUserID)

	}
}
