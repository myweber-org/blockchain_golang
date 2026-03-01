
package main

import (
	"fmt"
	"os"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.92,
			"GBP": 0.79,
			"JPY": 148.5,
			"CAD": 1.35,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromRate, fromExists := c.rates[fromCurrency]
	toRate, toExists := c.rates[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency: %s or %s", fromCurrency, toCurrency)
	}

	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for %s", fromCurrency)
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate
	return convertedAmount, nil
}

func (c *CurrencyConverter) AddRate(currency string, rate float64) error {
	if rate <= 0 {
		return fmt.Errorf("exchange rate must be positive")
	}
	c.rates[currency] = rate
	return nil
}

func (c *CurrencyConverter) ListCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
		fmt.Println("Available currencies:", converter.ListCurrencies())
		os.Exit(1)
	}

	var amount float64
	_, err := fmt.Sscanf(os.Args[1], "%f", &amount)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := converter.Convert(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}