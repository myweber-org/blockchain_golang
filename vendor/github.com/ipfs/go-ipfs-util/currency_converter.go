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
	fromRate, ok1 := er.rates[from]
	toRate, ok2 := er.rates[to]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unsupported currency")
	}

	if fromRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for source currency")
	}

	converted := (amount / fromRate) * toRate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) UpdateRate(currency Currency, rate float64) {
	if rate > 0 {
		er.rates[currency] = rate
	}
}

func main() {
	rates := NewExchangeRates()

	amount := 100.0
	result, err := rates.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	rates.UpdateRate(EUR, 0.88)
	newResult, _ := rates.Convert(amount, USD, EUR)
	fmt.Printf("After rate update: %.2f %s = %.2f %s\n", amount, USD, newResult, EUR)
}package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func fetchExchangeRates(apiKey string) (*ExchangeRates, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD")
	if apiKey != "" {
		url = fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", apiKey)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRates
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, err
	}

	return &rates, nil
}

func convertCurrency(amount float64, fromCurrency, toCurrency string, rates *ExchangeRates) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	fromRate, ok1 := rates.Rates[fromCurrency]
	toRate, ok2 := rates.Rates[toCurrency]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unsupported currency")
	}

	if rates.Base == fromCurrency {
		return amount * toRate, nil
	}

	amountInBase := amount / fromRate
	return amountInBase * toRate, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
		fmt.Println("Example: go run currency_converter.go 100 USD EUR")
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	rates, err := fetchExchangeRates(apiKey)
	if err != nil {
		fmt.Printf("Failed to fetch exchange rates: %v\n", err)
		os.Exit(1)
	}

	converted, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s (as of %s)\n", amount, fromCurrency, converted, toCurrency, rates.Date)
}