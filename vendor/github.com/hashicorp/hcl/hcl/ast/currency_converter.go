package main

import (
	"fmt"
	"sync"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates map[string]float64
	mu    sync.RWMutex
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := fmt.Sprintf("%s:%s", from, to)
	c.rates[key] = rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", from, to)
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) ListRates() []ExchangeRate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var rates []ExchangeRate
	for key, rate := range c.rates {
		var from, to string
		fmt.Sscanf(key, "%s:%s", &from, &to)
		rates = append(rates, ExchangeRate{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         rate,
		})
	}
	return rates
}

func main() {
	converter := NewCurrencyConverter()

	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.0)

	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)

	rates := converter.ListRates()
	fmt.Println("Available exchange rates:")
	for _, rate := range rates {
		fmt.Printf("%s to %s: %.4f\n", rate.FromCurrency, rate.ToCurrency, rate.Rate)
	}
}