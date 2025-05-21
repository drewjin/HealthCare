package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAddUserDataRouter(r *gin.Engine) {
	adduserdata := r.Group("/api/adduserdata")
	adduserdata.Use(middlewares.AuthMiddleWare())
	{
		// 机构用户或者管理员为用户添加体检数据
		adduserdata.POST("/:customer_id/:plan_id", middlewares.RequireUserType(3), controllers.AddUserDateByPlanID)
	}
}
