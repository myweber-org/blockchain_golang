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

type ExchangeRates struct {
	rates map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	return &ExchangeRates{
		rates: map[Currency]float64{
			USD: 1.0,
			EUR: 0.85,
			GBP: 0.73,
			JPY: 110.0,
		},
	}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	fromRate, okFrom := er.rates[from]
	toRate, okTo := er.rates[to]

	if !okFrom || !okTo {
		return 0, fmt.Errorf("unsupported currency")
	}

	baseAmount := amount / fromRate
	converted := baseAmount * toRate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) UpdateRate(currency Currency, rate float64) {
	if rate > 0 {
		er.rates[currency] = rate
	}
}

func main() {
	rates := NewExchangeRates()

	amount := 100.0
	converted, err := rates.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, EUR)

	rates.UpdateRate(EUR, 0.88)
	converted, _ = rates.Convert(amount, USD, EUR)
	fmt.Printf("After rate update: %.2f %s = %.2f %s\n", amount, USD, converted, EUR)
}