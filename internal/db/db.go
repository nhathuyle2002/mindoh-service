package db

import (
	"log"
	"mindoh-service/config"
	"mindoh-service/internal/expense"
	"mindoh-service/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	dsn := cfg.GetDSN()
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = database

	// Auto-migrate User model
	err = DB.AutoMigrate(&user.User{})
	err = DB.AutoMigrate(&expense.Expense{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
