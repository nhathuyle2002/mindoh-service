package currency

import (
	"mindoh-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterCurrencyRoutes(r *gin.Engine, a auth.IAuthService, resolveUser func(string) (uint, error)) {
	handler := NewCurrencyHandler()
	group := r.Group("/api/currency")
	group.Use(a.AuthMiddleware(resolveUser))
	{
		group.GET("/exchange-rates", handler.GetExchangeRates)
		group.GET("/currencies", handler.GetAvailableCurrencies)
	}
}
