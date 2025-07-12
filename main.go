package main

import (
	"mindoh-service/config"
	"mindoh-service/internal/db"
	"mindoh-service/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file before anything else
	godotenv.Load()

	cfg := config.LoadConfig()
	db.ConnectDatabase(cfg)

	dbInstance := db.GetDB()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	user.RegisterUserRoutes(r, dbInstance)

	r.Run()
}
