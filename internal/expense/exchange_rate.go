package expense

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
	Base  string             `json:"base"`  // For fallback API format
	Rates map[string]float64 `json:"rates"` // For fallback API format
	// For primary API format (nested structure)
	VND map[string]float64 `json:"vnd"`
}

var globalExchangeRateService *ExchangeRateService
var once sync.Once

// GetExchangeRateService returns the singleton instance
func GetExchangeRateService() *ExchangeRateService {
	once.Do(func() {
		globalExchangeRateService = &ExchangeRateService{
			rates:      make(map[string]float64),
			cacheTTL:   6 * time.Hour, // Cache for 6 hours
			lastUpdate: time.Time{},
		}
		// Initialize with default rates in case API fails
		globalExchangeRateService.rates = map[string]float64{
			"VND": 1,
			"USD": 25000,
			"EUR": 27000,
		}
		// Fetch initial rates
		go globalExchangeRateService.fetchRates()
	})
	return globalExchangeRateService
}

// fetchRates fetches exchange rates from the API
func (e *ExchangeRateService) fetchRates() {
	// Using fawazahmed0 currency API - free with no rate limits
	// Primary URL via CDN
	url := "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/vnd.json"

	resp, err := http.Get(url)
	if err != nil {
		// Try fallback URL
		url = "https://latest.currency-api.pages.dev/v1/currencies/vnd.json"
		resp, err = http.Get(url)
		if err != nil {
			fmt.Printf("Failed to fetch exchange rates from both URLs: %v\n", err)
			return
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Exchange rate API returned status: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read exchange rate response: %v\n", err)
		return
	}

	var apiResp ExchangeRateAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Printf("Failed to parse exchange rate response: %v\n", err)
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// The API returns rates FROM VND to other currencies
	// For example: 1 VND = 0.00004 USD means 1 USD = 25000 VND
	// We need to invert to get VND per currency
	newRates := make(map[string]float64)
	newRates["VND"] = 1 // VND to VND is always 1

	// The response has vnd object with currency rates
	if apiResp.VND != nil {
		if usdRate, ok := apiResp.VND["usd"]; ok && usdRate > 0 {
			newRates["USD"] = 1 / usdRate // Invert to get VND per USD
		}
		if eurRate, ok := apiResp.VND["eur"]; ok && eurRate > 0 {
			newRates["EUR"] = 1 / eurRate // Invert to get VND per EUR
		}
	}

	e.rates = newRates
	e.lastUpdate = time.Now()

	fmt.Printf("Exchange rates updated at %s: USD=%.2f VND, EUR=%.2f VND\n",
		e.lastUpdate.Format(time.RFC3339), newRates["USD"], newRates["EUR"])
}

// GetRates returns the current exchange rates (to VND)
// Fetches fresh rates if cache is older than 6 hours
func (e *ExchangeRateService) GetRates() map[string]float64 {
	e.mu.RLock()
	cacheExpired := time.Since(e.lastUpdate) > e.cacheTTL
	e.mu.RUnlock()

	// Only fetch if cache expired
	if cacheExpired {
		e.fetchRates()
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	// Return a copy to avoid concurrent access issues
	ratesCopy := make(map[string]float64)
	for k, v := range e.rates {
		ratesCopy[k] = v
	}
	return ratesCopy
}

// RefreshRates forces a refresh of exchange rates
func (e *ExchangeRateService) RefreshRates() {
	e.fetchRates()
}
