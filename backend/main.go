package main

import (
	"fmt"
	"healthcare/config"
	"healthcare/routers"
)

func main() {
	config.InitConfig()
	fmt.Printf("Loaded config: %+v\n\n", config.AppConfig)
	router := routers.SetupRouter()
	port := config.AppConfig.App.BackendPort
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
