package user

import (
	"mindoh-service/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var authService *AuthService

type AuthContext struct {
	UserID uint `json:"user_id"`
	Role   Role `json:"role"`
}

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
		// Set AuthContext in context
		authCtx := AuthContext{
			UserID: userID,
			Role:   role,
		}
		c.Set("auth", authCtx)
		c.Next()
	}
}

// GetAuthContext retrieves the AuthContext from Gin context. Returns zero AuthContext if not set or wrong type.
func GetAuthContext(c *gin.Context) AuthContext {
	val, exists := c.Get("auth")
	if !exists {
		return AuthContext{}
	}
	authCtx, ok := val.(AuthContext)
	if !ok {
		return AuthContext{}
	}
	return authCtx
}
