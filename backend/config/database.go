package config

import (
	"HealthCare/backend/global"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	dsn := AppConfig.Database.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database, got error: %v", err)
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Fatalf("Failed to configure database, got error: %v", err)
	}

	// // 自动迁移数据库表
	// err = db.AutoMigrate(&models.User{}, &models.Family{}, &models.Institution{}, &models.Plan{}, &models.HealthItem{}, &models.PlanHeathItem{}, &models.Commentary{}, &models.UserHealthItem{}, &models.UserPackage{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database tables, got error: %v", err)
	// }

	global.DB = db
}
