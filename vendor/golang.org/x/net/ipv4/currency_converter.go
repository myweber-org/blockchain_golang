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
	fromRate, fromOk := er.rates[from]
	toRate, toOk := er.rates[to]

	if !fromOk || !toOk {
		return 0, fmt.Errorf("unsupported currency")
	}

	baseAmount := amount / fromRate
	converted := baseAmount * toRate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) AddRate(currency Currency, rate float64) {
	er.rates[currency] = rate
}

func main() {
	rates := NewExchangeRates()

	amount := 100.0
	result, err := rates.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	rates.AddRate("CAD", 1.25)
	cadResult, _ := rates.Convert(50.0, USD, "CAD")
	fmt.Printf("50.00 %s = %.2f CAD\n", USD, cadResult)
}