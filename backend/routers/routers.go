package routers

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	routers := gin.Default()
	SetupAuthRouter(routers)
	SetupUserRouter(routers)
	return routers
}
