package user

import (
	"mindoh-service/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB, authService IAuthService) {
	handler := NewUserHandler(authService, db)

	// Public routes
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)

	// Protected routes
	auth := r.Group("/api")
	auth.Use(AuthMiddleware(config.LoadConfig()))
	{
		auth.GET("/users/:id", handler.GetUser)
		auth.PUT("/users/:id", handler.UpdateUser)
		auth.DELETE("/users/:id", handler.DeleteUser)
	}
}
