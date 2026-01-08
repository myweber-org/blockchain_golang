
package main

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

	fromRate, fromExists := rates.Rates[fromCurrency]
	toRate, toExists := rates.Rates[toCurrency]

	if !fromExists || !toExists {
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
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
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
		fmt.Println("Using fallback rates...")
		rates = &ExchangeRates{
			Base: "USD",
			Rates: map[string]float64{
				"USD": 1.0,
				"EUR": 0.85,
				"GBP": 0.73,
				"JPY": 110.0,
				"CAD": 1.25,
			},
		}
	}

	result, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
	fmt.Printf("Based on exchange rates from %s\n", rates.Date)
}