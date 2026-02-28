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
			"EUR": 0.85,
			"GBP": 0.73,
			"JPY": 110.0,
			"CAD": 1.25,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	fromRate, fromOk := c.rates[fromCurrency]
	toRate, toOk := c.rates[toCurrency]

	if !fromOk {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}
	if !toOk {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for currency: %s", fromCurrency)
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
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

	// Example conversion
	result, err := converter.Convert(100.0, "USD", "EUR")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("100 USD = %.2f EUR\n", result)
	fmt.Printf("Available currencies: %v\n", converter.ListCurrencies())

	// Add new currency
	err = converter.AddRate("AUD", 1.35)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add rate: %v\n", err)
	}

	// Convert using new currency
	result, err = converter.Convert(50.0, "AUD", "JPY")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Conversion error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("50 AUD = %.2f JPY\n", result)
}