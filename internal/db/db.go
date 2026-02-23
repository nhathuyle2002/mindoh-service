package db

import (
	"log/slog"
	"mindoh-service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	dsn := cfg.GetDSN()
	slog.Info("connecting to database", "host", cfg.DB.Host, "port", cfg.DB.Port, "name", cfg.DB.Name)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		panic("database connection failed")
	}
	DB = database
	slog.Info("database connected")

	// Auto-migrate models
	if err := DB.AutoMigrate(&User{}, &Expense{}); err != nil {
		slog.Error("failed to migrate database", "error", err)
		panic("database migration failed")
	}
	slog.Info("database migration ok")
}

func GetDB() *gorm.DB {
	return DB
}
