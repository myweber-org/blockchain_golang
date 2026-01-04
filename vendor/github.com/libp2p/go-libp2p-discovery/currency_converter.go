package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiEndpoint string
	client      *http.Client
	cache       map[string]float64
	lastUpdated time.Time
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		apiEndpoint: "https://api.exchangerate-api.com/v4/latest/USD",
		client:      &http.Client{Timeout: 10 * time.Second},
		cache:       make(map[string]float64),
	}
}

func (c *CurrencyConverter) fetchRates() error {
	resp, err := c.client.Get(c.apiEndpoint)
	if err != nil {
		return fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var rateResponse ExchangeRateResponse
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	c.cache = rateResponse.Rates
	c.cache[rateResponse.Base] = 1.0
	c.lastUpdated = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if time.Since(c.lastUpdated) > 30*time.Minute || len(c.cache) == 0 {
		if err := c.fetchRates(); err != nil {
			return 0, err
		}
	}

	fromRate, fromExists := c.cache[fromCurrency]
	toRate, toExists := c.cache[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("invalid currency code: %s or %s", fromCurrency, toCurrency)
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}