package currency

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CurrencyHandler struct{}

func NewCurrencyHandler() *CurrencyHandler { return &CurrencyHandler{} }

// GetExchangeRates godoc
// @Summary Get exchange rates
// @Description Get exchange rates (base VND)
// @Tags currency
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /currency/exchange-rates [get]
func (h *CurrencyHandler) GetExchangeRates(c *gin.Context) {
	rates := GetExchangeRateService().GetRates()
	c.JSON(http.StatusOK, gin.H{"base_currency": "VND", "rates": rates})
}

// GetAvailableCurrencies godoc
// @Summary Get available currencies
// @Description List supported currency codes
// @Tags currency
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /currency/currencies [get]
func (h *CurrencyHandler) GetAvailableCurrencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"currencies": AvailableCurrencies})
}
