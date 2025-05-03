package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Author  string `yaml:"author"`
		Port    string `yaml:"port"`
	}

	Database struct {
		DSN          string `yaml:"dsn"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	}

	Relations []string

	EMail struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		ServerEmail string `yaml:"server_email"`
		Password    string `yaml:"password"`
		From        string `yaml:"from"`
		FromName    string `yaml:"from_name"`
	}
}

var (
	AppConfig      *Config
	ValidRelations = make(map[string]bool)
)

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

	for _, rel := range AppConfig.Relations {
		ValidRelations[rel] = true
	}

	InitDB()
}
