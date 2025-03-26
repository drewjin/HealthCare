package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name    string
		Version string
		Author  string
		Port    string
	}
	Database struct {
		DSN          string
		MaxIdleConns int
		MaxOpenConns int
	}
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	InitDB()
}
