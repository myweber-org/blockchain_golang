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
			"USD": 1.0,
			"EUR": 0.85,
			"GBP": 0.73,
			"JPY": 110.0,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	fromRate, ok1 := c.rates[from]
	toRate, ok2 := c.rates[to]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unsupported currency")
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) {
	c.rates[currency] = rate
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	converter.AddRate("CAD", 1.25)
	result, _ = converter.Convert(50.0, "CAD", "JPY")
	fmt.Printf("50.00 CAD = %.2f JPY\n", result)
}