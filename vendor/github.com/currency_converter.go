package main

import (
	"fmt"
	"math"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type ExchangeRates map[Currency]float64

type CurrencyConverter struct {
	rates ExchangeRates
}

func NewCurrencyConverter(rates ExchangeRates) *CurrencyConverter {
	return &CurrencyConverter{rates: rates}
}

func (c *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	fromRate, ok := c.rates[from]
	if !ok {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}

	toRate, ok := c.rates[to]
	if !ok {
		return 0, fmt.Errorf("unsupported target currency: %s", to)
	}

	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for %s", from)
	}

	converted := (amount / fromRate) * toRate
	return math.Round(converted*100) / 100, nil
}

func (c *CurrencyConverter) UpdateRates(newRates ExchangeRates) {
	for currency, rate := range newRates {
		c.rates[currency] = rate
	}
}

func main() {
	rates := ExchangeRates{
		USD: 1.0,
		EUR: 0.85,
		GBP: 0.73,
		JPY: 110.5,
	}

	converter := NewCurrencyConverter(rates)

	amount := 100.0
	result, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	newRates := ExchangeRates{EUR: 0.88, GBP: 0.75}
	converter.UpdateRates(newRates)

	result, _ = converter.Convert(amount, USD, EUR)
	fmt.Printf("After rate update: %.2f %s = %.2f %s\n", amount, USD, result, EUR)
}