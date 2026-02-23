package user

import (
	"mindoh-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, authService auth.IAuthService, userService *UserService, resolveUser func(string) (uint, error)) {
	handler := NewUserHandler(authService, userService)

	// Public routes
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)
	r.GET("/api/verify-email", handler.VerifyEmail)
	r.POST("/api/resend-verification", handler.ResendVerification)
	r.POST("/api/forgot-password", handler.ForgotPassword)
	r.POST("/api/reset-password", handler.ResetPassword)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authService.AuthMiddleware(resolveUser))
	{
		protected.GET("/users/me", handler.GetMe)
		protected.PUT("/users/me", handler.UpdateMe)
		protected.GET("/users/:id", handler.GetUser)
		protected.PUT("/users/:id", handler.UpdateUser)
		protected.DELETE("/users/:id", handler.DeleteUser)
		protected.POST("/users/change-password", handler.ChangePassword)
	}

	// Admin-only routes
	admin := r.Group("/api/admin")
	admin.Use(authService.AuthMiddleware(resolveUser), authService.RoleGuard(auth.RoleAdmin))
	{
		admin.POST("/users", handler.AdminCreateUser)
	}
}
