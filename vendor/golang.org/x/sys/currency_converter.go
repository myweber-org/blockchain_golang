
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	var usdAmount float64
	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}package main

import (
	"fmt"
)

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD_EUR": 0.85,
			"EUR_USD": 1.18,
			"USD_GBP": 0.73,
			"GBP_USD": 1.37,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("conversion rate not available for %s to %s", from, to)
	}
	return amount * rate, nil
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.rates[from+"_"+to] = rate
}

func main() {
	converter := NewCurrencyConverter()
	
	converted, err := converter.Convert(100, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("100 USD = %.2f EUR\n", converted)
	
	converter.AddRate("EUR_GBP", 0.86)
	
	converted, err = converter.Convert(50, "EUR", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("50 EUR = %.2f GBP\n", converted)
}