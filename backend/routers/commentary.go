package routers

import (
	"healthcare/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCommentaryRouter(r *gin.Engine) {
	commentary := r.Group("/api/commentary")
	commentary.Use(middlewares.AuthMiddleWare())
	{
		// 发布评论

		// 删除评论
	}
}
