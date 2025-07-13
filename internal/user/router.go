package user

import (
	"mindoh-service/config"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, cfg *config.Config, authService IAuthService, userService *UserService) {
	handler := NewUserHandler(authService, userService)

	// Public routes
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)

	// Protected routes
	auth := r.Group("/api")
	auth.Use(AuthMiddleware(cfg))
	{
		auth.GET("/users/:id", handler.GetUser)
		auth.PUT("/users/:id", handler.UpdateUser)
		auth.DELETE("/users/:id", handler.DeleteUser)
	}
}
