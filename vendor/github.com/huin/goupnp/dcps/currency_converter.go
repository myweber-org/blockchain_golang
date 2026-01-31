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
			"JPY": 110.0,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	fromRate, ok := c.rates[from]
	if !ok {
		return 0, fmt.Errorf("unknown currency: %s", from)
	}
	toRate, ok := c.rates[to]
	if !ok {
		return 0, fmt.Errorf("unknown currency: %s", to)
	}
	return amount * (toRate / fromRate), nil
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) {
	c.rates[currency] = rate
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("CAD", 1.25)
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	result, _ = converter.Convert(amount, "EUR", "GBP")
	fmt.Printf("%.2f EUR = %.2f GBP\n", amount, result)
}
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}