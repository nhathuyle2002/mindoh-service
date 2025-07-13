package user

import (
	"mindoh-service/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var authService *AuthService

// AuthMiddleware checks JWT authentication and sets user info in context
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Initialize auth service with the provided config
	authService = NewAuthService(cfg)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := authService.ParseAndValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		// Set user info in context
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
