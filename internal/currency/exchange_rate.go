package currency

import (
	"encoding/json"
	"fmt"
	"io"
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
	url := "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/vnd.json"
	resp, err := http.Get(url)
	if err != nil {
		url = "https://latest.currency-api.pages.dev/v1/currencies/vnd.json"
		resp, err = http.Get(url)
		if err != nil {
			fmt.Printf("Failed to fetch exchange rates: %v\n", err)
			return
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Exchange rate API status: %d\n", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed reading exchange rate body: %v\n", err)
		return
	}
	var apiResp ExchangeRateAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Printf("Failed parsing exchange rate JSON: %v\n", err)
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
	fmt.Printf("Exchange rates updated: USD=%.2f VND EUR=%.2f VND\n", newRates["USD"], newRates["EUR"])
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
