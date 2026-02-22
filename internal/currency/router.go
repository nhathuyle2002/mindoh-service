package currency

import (
	"mindoh-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterCurrencyRoutes(r *gin.Engine, a auth.IAuthService) {
	handler := NewCurrencyHandler()
	group := r.Group("/api/currency")
	group.Use(a.AuthMiddleware())
	{
		group.GET("/exchange-rates", handler.GetExchangeRates)
		group.GET("/currencies", handler.GetAvailableCurrencies)
	}
}
