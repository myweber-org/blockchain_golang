package main

import (
	"fmt"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

var exchangeRates = map[Currency]float64{
	USD: 1.0,
	EUR: 0.85,
	GBP: 0.73,
}

func Convert(amount float64, from, to Currency) (float64, error) {
	fromRate, ok := exchangeRates[from]
	if !ok {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}
	toRate, ok := exchangeRates[to]
	if !ok {
		return 0, fmt.Errorf("unsupported target currency: %s", to)
	}
	return amount * toRate / fromRate, nil
}

func main() {
	amount := 100.0
	result, err := Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)
}