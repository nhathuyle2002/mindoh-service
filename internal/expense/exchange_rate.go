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

// ExchangeRateAPIResponse represents the response from exchangerate-api.com
type ExchangeRateAPIResponse struct {
	Result            string             `json:"result"`
	BaseCode          string             `json:"base_code"`
	ConversionRates   map[string]float64 `json:"conversion_rates"`
	TimeLastUpdateUTC string             `json:"time_last_update_utc"`
}

var globalExchangeRateService *ExchangeRateService
var once sync.Once

// GetExchangeRateService returns the singleton instance
func GetExchangeRateService() *ExchangeRateService {
	once.Do(func() {
		globalExchangeRateService = &ExchangeRateService{
			rates:      make(map[string]float64),
			cacheTTL:   1 * time.Hour, // Cache for 1 hour
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
	// Using exchangerate-api.com free tier (1500 requests/month)
	// Base currency is VND
	url := "https://api.exchangerate-api.com/v4/latest/VND"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch exchange rates: %v\n", err)
		return
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

	if apiResp.Result != "success" {
		fmt.Printf("Exchange rate API returned non-success result: %s\n", apiResp.Result)
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Convert rates to VND base
	// The API returns rates FROM VND, so we need to invert them
	// For example: 1 VND = 0.00004 USD means 1 USD = 25000 VND
	newRates := make(map[string]float64)
	newRates["VND"] = 1 // VND to VND is always 1

	if usdRate, ok := apiResp.ConversionRates["USD"]; ok && usdRate > 0 {
		newRates["USD"] = 1 / usdRate // Invert to get VND per USD
	}
	if eurRate, ok := apiResp.ConversionRates["EUR"]; ok && eurRate > 0 {
		newRates["EUR"] = 1 / eurRate // Invert to get VND per EUR
	}

	e.rates = newRates
	e.lastUpdate = time.Now()

	fmt.Printf("Exchange rates updated at %s: USD=%.2f VND, EUR=%.2f VND\n",
		e.lastUpdate.Format(time.RFC3339), newRates["USD"], newRates["EUR"])
}

// GetRates returns the current exchange rates (to VND)
// Always fetches fresh rates from API
func (e *ExchangeRateService) GetRates() map[string]float64 {
	// Always fetch fresh rates when called
	e.fetchRates()

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
