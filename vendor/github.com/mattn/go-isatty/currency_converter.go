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
	cache       map[string]ExchangeRateResponse
	lastUpdated time.Time
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		apiEndpoint: "https://api.exchangerate-api.com/v4/latest/",
		client:      &http.Client{Timeout: 10 * time.Second},
		cache:       make(map[string]ExchangeRateResponse),
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	rates, err := c.getExchangeRates(fromCurrency)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not supported", toCurrency)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) getExchangeRates(baseCurrency string) (*ExchangeRateResponse, error) {
	if cached, exists := c.cache[baseCurrency]; exists {
		if time.Since(c.lastUpdated) < 30*time.Minute {
			return &cached, nil
		}
	}

	resp, err := c.client.Get(c.apiEndpoint + baseCurrency)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRateResponse
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, err
	}

	c.cache[baseCurrency] = rates
	c.lastUpdated = time.Now()

	return &rates, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies(baseCurrency string) ([]string, error) {
	rates, err := c.getExchangeRates(baseCurrency)
	if err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(rates.Rates))
	for currency := range rates.Rates {
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	currencies, err := converter.GetSupportedCurrencies("USD")
	if err != nil {
		fmt.Printf("Error fetching currencies: %v\n", err)
		return
	}

	fmt.Printf("Supported currencies based on USD: %v\n", currencies[:10])
}