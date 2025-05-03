package middlewares

import (
	"fmt"
	"healthcare/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCorsMiddleware() gin.HandlerFunc {
	port := config.AppConfig.App.FrontendPort
	if port == "" {
		port = "5173"
	}
	URL := fmt.Sprintf("http://localhost:%s", port)
	return cors.New(cors.Config{
		AllowOrigins:     []string{URL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
