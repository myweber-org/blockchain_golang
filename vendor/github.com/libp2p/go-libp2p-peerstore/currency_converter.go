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

type Converter struct {
	rates map[Currency]map[Currency]float64
}

func NewConverter() *Converter {
	c := &Converter{
		rates: make(map[Currency]map[Currency]float64),
	}
	
	c.rates[USD] = map[Currency]float64{
		EUR: 0.92,
		GBP: 0.79,
	}
	
	c.rates[EUR] = map[Currency]float64{
		USD: 1.09,
		GBP: 0.86,
	}
	
	c.rates[GBP] = map[Currency]float64{
		USD: 1.27,
		EUR: 1.16,
	}
	
	return c
}

func (c *Converter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}
	
	rateMap, exists := c.rates[from]
	if !exists {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}
	
	rate, exists := rateMap[to]
	if !exists {
		return 0, fmt.Errorf("conversion from %s to %s not supported", from, to)
	}
	
	return amount * rate, nil
}

func main() {
	converter := NewConverter()
	
	amount := 100.0
	
	usdToEur, _ := converter.Convert(amount, USD, EUR)
	usdToGbp, _ := converter.Convert(amount, USD, GBP)
	eurToGbp, _ := converter.Convert(amount, EUR, GBP)
	
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, usdToEur, EUR)
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, usdToGbp, GBP)
	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, eurToGbp, GBP)
}