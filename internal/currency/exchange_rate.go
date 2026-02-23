package currency

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// ExchangeRateService handles fetching and caching exchange rates
type ExchangeRateService struct {
	rates      map[string]float64
	lastUpdate time.Time
	mu         sync.RWMutex
	cacheTTL   time.Duration
}

// ExchangeRateAPIResponse represents the response from fawazahmed0 currency API
type ExchangeRateAPIResponse struct {
	Date  string             `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	VND   map[string]float64 `json:"vnd"`
}

var globalExchangeRateService *ExchangeRateService
var once sync.Once

// GetExchangeRateService returns the singleton instance
func GetExchangeRateService() *ExchangeRateService {
	once.Do(func() {
		globalExchangeRateService = &ExchangeRateService{
			rates:      make(map[string]float64),
			cacheTTL:   6 * time.Hour,
			lastUpdate: time.Time{},
		}
		globalExchangeRateService.rates = map[string]float64{
			"VND": 1,
			"USD": 25000,
			"EUR": 27000,
		}
		go globalExchangeRateService.fetchRates()
	})
	return globalExchangeRateService
}

func (e *ExchangeRateService) fetchRates() {
	primaryURL := "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/vnd.json"
	fallbackURL := "https://latest.currency-api.pages.dev/v1/currencies/vnd.json"

	slog.Info("fetching exchange rates", "url", primaryURL)
	resp, err := http.Get(primaryURL)
	if err != nil {
		slog.Warn("primary exchange rate URL failed, trying fallback", "error", err, "url", fallbackURL)
		resp, err = http.Get(fallbackURL)
		if err != nil {
			slog.Error("failed to fetch exchange rates", "error", err)
			return
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("exchange rate API returned non-200", "status", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed reading exchange rate body", "error", err)
		return
	}
	var apiResp ExchangeRateAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		slog.Error("failed parsing exchange rate JSON", "error", err)
		return
	}
	newRates := make(map[string]float64)
	newRates["VND"] = 1
	if apiResp.VND != nil {
		if usdRate, ok := apiResp.VND["usd"]; ok && usdRate > 0 {
			newRates["USD"] = 1 / usdRate
		}
		if eurRate, ok := apiResp.VND["eur"]; ok && eurRate > 0 {
			newRates["EUR"] = 1 / eurRate
		}
	}
	e.mu.Lock()
	e.rates = newRates
	e.lastUpdate = time.Now()
	e.mu.Unlock()
	slog.Info("exchange rates updated", "USD_to_VND", newRates["USD"], "EUR_to_VND", newRates["EUR"])
}

// GetRates returns a copy of current rates (to VND)
func (e *ExchangeRateService) GetRates() map[string]float64 {
	e.mu.RLock()
	expired := time.Since(e.lastUpdate) > e.cacheTTL
	e.mu.RUnlock()
	if expired {
		e.fetchRates()
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make(map[string]float64)
	for k, v := range e.rates {
		out[k] = v
	}
	return out
}

// RefreshRates forces refresh
func (e *ExchangeRateService) RefreshRates() { e.fetchRates() }
