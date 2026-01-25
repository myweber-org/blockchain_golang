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
	rates map[Currency]map[Currency]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[Currency]map[Currency]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to Currency, rate float64) {
	if c.rates[from] == nil {
		c.rates[from] = make(map[Currency]float64)
	}
	c.rates[from][to] = rate

	// Add inverse rate
	if c.rates[to] == nil {
		c.rates[to] = make(map[Currency]float64)
	}
	c.rates[to][from] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	if rate, exists := c.rates[from][to]; exists {
		return round(amount*rate, 2), nil
	}

	return 0, fmt.Errorf("no exchange rate found from %s to %s", from, to)
}

func round(value float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(value*multiplier) / multiplier
}

func main() {
	converter := NewCurrencyConverter()

	// Add sample exchange rates
	converter.AddRate(USD, EUR, 0.92)
	converter.AddRate(USD, GBP, 0.79)
	converter.AddRate(USD, JPY, 149.5)
	converter.AddRate(EUR, GBP, 0.86)

	// Perform conversions
	amounts := []float64{100, 250, 500}
	conversions := []struct {
		from Currency
		to   Currency
	}{
		{USD, EUR},
		{USD, GBP},
		{EUR, GBP},
		{GBP, JPY},
	}

	for _, amount := range amounts {
		fmt.Printf("Converting %.2f:\n", amount)
		for _, conv := range conversions {
			result, err := converter.Convert(amount, conv.from, conv.to)
			if err != nil {
				fmt.Printf("  %s to %s: %v\n", conv.from, conv.to, err)
			} else {
				fmt.Printf("  %s to %s: %.2f\n", conv.from, conv.to, result)
			}
		}
		fmt.Println()
	}
}