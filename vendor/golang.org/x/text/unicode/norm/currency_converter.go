package main

import (
	"encoding/json"
	"fmt"
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
	client    *http.Client
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		client: &http.Client{Timeout: 10 * time.Second},
		rates:  make(map[string]float64),
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < 30*time.Minute && len(c.rates) > 0 {
		return nil
	}

	resp, err := c.client.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var rates ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return err
	}

	c.rates = rates.Rates
	c.rates["USD"] = 1.0
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, ok1 := c.rates[from]
	toRate, ok2 := c.rates[to]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid currency code")
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
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
}