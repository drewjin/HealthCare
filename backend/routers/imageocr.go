package routers

import (
	"healthcare/controllers"
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupImageOcrRouter(routers *gin.Engine) {
	imageocr := routers.Group("/api/imageocr")
	imageocr.Use(middlewares.AuthMiddleWare())
	{
		imageocr.POST("/solve", controllers.ImageOcr)
	}
}
