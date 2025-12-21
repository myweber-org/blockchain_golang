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
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	fromRate, ok := c.rates[from]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", from)
	}
	toRate, ok := c.rates[to]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", to)
	}
	return amount * (toRate / fromRate), nil
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
	
	result, err = converter.Convert(amount, "USD", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f GBP\n", amount, result)
}