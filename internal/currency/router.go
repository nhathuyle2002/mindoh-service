package currency

import (
	"mindoh-service/internal/auth"

	"github.com/gin-gonic/gin"
)

// BypassOrAuth allows requests with header `X-By-Pass` to skip auth; otherwise enforces JWT auth.
func BypassOrAuth(a auth.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-By-Pass") != "" {
			// Skip auth and continue
			c.Next()
			return
		}
		// Enforce auth when bypass header absent
		a.AuthMiddleware()(c)
	}
}

func RegisterCurrencyRoutes(r *gin.Engine, a auth.IAuthService) {
	handler := NewCurrencyHandler()
	group := r.Group("/api/currency")
	group.Use(BypassOrAuth(a))
	{
		group.GET("/exchange-rates", handler.GetExchangeRates)
		group.GET("/currencies", handler.GetAvailableCurrencies)
	}
}
