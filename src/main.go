package main

import (
	"healthcare/config"
	"healthcare/routers"
)

func main() {
	config.InitConfig()
	router := routers.SetupRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}
	router.Run(port)
}
