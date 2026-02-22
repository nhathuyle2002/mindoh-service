package auth

import (
	"mindoh-service/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	cfg *config.Config
}

type IAuthService interface {
	AuthMiddleware() gin.HandlerFunc
	RoleGuard(roles ...Role) gin.HandlerFunc
	GenerateJWT(userID uint, role Role) (string, error)
	ParseAndValidateJWT(tokenString string) (uint, Role, error)
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

// AuthMiddleware checks JWT authentication and sets user info in context
func (a *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := a.ParseAndValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
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

// RequireRole is a middleware that checks if the user has the required role
func (a *AuthService) RoleGuard(roles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx := GetAuthContext(c)
		if authCtx.Role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		for _, role := range roles {
			if authCtx.Role == role {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}
