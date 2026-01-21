
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	rates     map[string]float64
	lastFetch time.Time
	cacheTTL  time.Duration
}

func NewCurrencyConverter(ttl time.Duration) *CurrencyConverter {
	return &CurrencyConverter{
		rates:    make(map[string]float64),
		cacheTTL: ttl,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < c.cacheTTL && len(c.rates) > 0 {
		return nil
	}

	resp, err := http.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var exchangeRates ExchangeRates
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		return fmt.Errorf("failed to parse rates: %w", err)
	}

	c.rates = exchangeRates.Rates
	c.rates["USD"] = 1.0
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, ok := c.rates[from]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", from)
	}

	toRate, ok := c.rates[to]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", to)
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	if err := c.fetchRates(); err != nil {
		return []string{}
	}

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter(30 * time.Minute)

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	fmt.Println("Supported currencies:")
	for _, currency := range converter.GetSupportedCurrencies() {
		fmt.Printf("- %s\n", currency)
	}
}