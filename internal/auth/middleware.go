package auth

import (
	"log/slog"
	"mindoh-service/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	cfg *config.Config
}

type IAuthService interface {
	AuthMiddleware(resolveUser func(username string) (uint, error)) gin.HandlerFunc
	RoleGuard(roles ...Role) gin.HandlerFunc
	GenerateJWT(username string, role Role) (string, error)
	ParseAndValidateJWT(tokenString string) (string, Role, error)
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

// AuthMiddleware checks JWT authentication and sets user info in context
func (a *AuthService) AuthMiddleware(resolveUser func(username string) (uint, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Warn("missing or invalid authorization header", "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		username, role, err := a.ParseAndValidateJWT(tokenString)
		if err != nil {
			slog.Warn("invalid or expired JWT", "path", c.Request.URL.Path, "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		userID, err := resolveUser(username)
		if err != nil {
			slog.Warn("JWT username not found in DB", "username", username, "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		authCtx := AuthContext{
			UserID:   userID,
			Username: username,
			Role:     role,
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
