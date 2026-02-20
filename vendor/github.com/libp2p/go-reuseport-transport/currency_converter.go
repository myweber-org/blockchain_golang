
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
	fromRate, ok := er.rates[from]
	if !ok {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}
	
	toRate, ok := er.rates[to]
	if !ok {
		return 0, fmt.Errorf("unsupported target currency: %s", to)
	}
	
	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for currency: %s", from)
	}
	
	converted := (amount / fromRate) * toRate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) UpdateRate(currency Currency, rate float64) error {
	if rate <= 0 {
		return fmt.Errorf("exchange rate must be positive")
	}
	er.rates[currency] = rate
	return nil
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
	
	err = rates.UpdateRate(JPY, 115.5)
	if err != nil {
		fmt.Printf("Update error: %v\n", err)
	}
}