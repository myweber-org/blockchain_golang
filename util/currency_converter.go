
package main

import (
	"fmt"
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
		rates: make(map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) {
	c.rates[currency] = rate
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	fromRate, fromExists := c.rates[fromCurrency]
	toRate, toExists := c.rates[toCurrency]

	if !fromExists {
		return 0, fmt.Errorf("exchange rate for %s not found", fromCurrency)
	}
	if !toExists {
		return 0, fmt.Errorf("exchange rate for %s not found", toCurrency)
	}

	baseAmount := amount / fromRate
	return baseAmount * toRate, nil
}

func main() {
	converter := NewCurrencyConverter()
	converter.AddRate("USD", 1.0)
	converter.AddRate("EUR", 0.85)
	converter.AddRate("GBP", 0.75)
	converter.AddRate("JPY", 110.0)

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	result, err = converter.Convert(amount, "EUR", "GBP")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f GBP\n", amount, result)
}