package main

import (
	"fmt"
)

type CurrencyConverter struct {
	exchangeRates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		exchangeRates: map[string]float64{
			"USD_EUR": 0.92,
			"EUR_USD": 1.09,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	key := from + "_" + to
	rate, exists := c.exchangeRates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}
	return amount * rate, nil
}

func main() {
	converter := NewCurrencyConverter()
	
	usdAmount := 100.0
	converted, err := converter.Convert(usdAmount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, converted)
}