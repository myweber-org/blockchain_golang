
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

type ExchangeRate struct {
	From Currency
	To   Currency
	Rate float64
}

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	cc := &CurrencyConverter{
		rates: make(map[string]float64),
	}
	cc.initializeRates()
	return cc
}

func (cc *CurrencyConverter) initializeRates() {
	baseRates := []ExchangeRate{
		{USD, EUR, 0.85},
		{USD, GBP, 0.73},
		{USD, JPY, 110.5},
		{EUR, USD, 1.18},
		{EUR, GBP, 0.86},
		{EUR, JPY, 130.2},
		{GBP, USD, 1.37},
		{GBP, EUR, 1.16},
		{GBP, JPY, 151.8},
		{JPY, USD, 0.0091},
		{JPY, EUR, 0.0077},
		{JPY, GBP, 0.0066},
	}

	for _, rate := range baseRates {
		key := cc.getRateKey(rate.From, rate.To)
		cc.rates[key] = rate.Rate
	}
}

func (cc *CurrencyConverter) getRateKey(from, to Currency) string {
	return string(from) + "_" + string(to)
}

func (cc *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := cc.getRateKey(from, to)
	rate, exists := cc.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not available for %s to %s", from, to)
	}

	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (cc *CurrencyConverter) AddRate(from, to Currency, rate float64) {
	key := cc.getRateKey(from, to)
	cc.rates[key] = rate
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := USD
	to := EUR

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	converter.AddRate(USD, CAD, 1.25)
	cadResult, _ := converter.Convert(amount, USD, "CAD")
	fmt.Printf("%.2f %s = %.2f CAD\n", amount, from, cadResult)
}