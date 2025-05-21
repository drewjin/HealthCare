package routers

import (
	"HealthCare/backend/controllers"
	"HealthCare/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupImageOcrRouter(routers *gin.Engine) {
	imageocr := routers.Group("/api/imageocr")
	imageocr.Use(middlewares.AuthMiddleWare())
	{
		imageocr.POST("/solve", controllers.ImageOcr)
	}
}
