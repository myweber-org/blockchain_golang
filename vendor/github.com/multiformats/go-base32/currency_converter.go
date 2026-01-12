
package main

import (
	"fmt"
	"os"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

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
			"CAD": 1.25,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromRate, fromOk := c.rates[fromCurrency]
	toRate, toOk := c.rates[toCurrency]

	if !fromOk || !toOk {
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

	result, err := converter.Convert(100.0, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("100 USD = %.2f EUR\n", result)

	converter.AddRate("AUD", 1.35)
	audResult, _ := converter.Convert(50.0, "AUD", "JPY")
	fmt.Printf("50 AUD = %.2f JPY\n", audResult)
}