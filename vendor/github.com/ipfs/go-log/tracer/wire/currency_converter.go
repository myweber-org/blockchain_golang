package main

import (
	"fmt"
)

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD_EUR": 0.85,
			"EUR_USD": 1.18,
			"USD_GBP": 0.73,
			"GBP_USD": 1.37,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not available for %s to %s", from, to)
	}
	return amount * rate, nil
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.rates[from+"_"+to] = rate
}

func main() {
	converter := NewCurrencyConverter()
	
	// Convert 100 USD to EUR
	result, err := converter.Convert(100, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("100 USD = %.2f EUR\n", result)
	
	// Add custom rate and convert
	converter.AddRate("USD_JPY", 110.5)
	jpyResult, _ := converter.Convert(50, "USD", "JPY")
	fmt.Printf("50 USD = %.2f JPY\n", jpyResult)
}