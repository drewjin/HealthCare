package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name         string `mapstructure:"name" yaml:"name"`
		Version      string `mapstructure:"version" yaml:"version"`
		Author       string `mapstructure:"author" yaml:"author"`
		BackendPort  string `mapstructure:"backend-port" yaml:"backend-port"`
		FrontendPort string `mapstructure:"frontend-port" yaml:"frontend-port"`
	} `mapstructure:"app" yaml:"app"`

	Database struct {
		DSN          string `mapstructure:"dsn" yaml:"dsn"`
		MaxIdleConns int    `mapstructure:"max-idle-conns" yaml:"max-idle-conns"`
		MaxOpenConns int    `mapstructure:"max-open-conns" yaml:"max-open-conns"`
	} `mapstructure:"database" yaml:"database"`

	Relations []string `mapstructure:"relations" yaml:"relations"`

	EMail struct {
		Host        string `mapstructure:"host" yaml:"host"`
		Port        int    `mapstructure:"port" yaml:"port"`
		ServerEmail string `mapstructure:"server-email" yaml:"server-email"`
		Password    string `mapstructure:"password" yaml:"password"`
		From        string `mapstructure:"from" yaml:"from"`
		FromName    string `mapstructure:"from-name" yaml:"from-name"`
	} `mapstructure:"email" yaml:"email"`
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
