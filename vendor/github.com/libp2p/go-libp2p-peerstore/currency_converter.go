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
	url := fmt.Sprintf("https://api.exchangeratesapi.io/latest?access_key=%s", apiKey)
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
	if fromCurrency == rates.Base {
		rate, exists := rates.Rates[toCurrency]
		if !exists {
			return 0, fmt.Errorf("currency %s not found in exchange rates", toCurrency)
		}
		return amount * rate, nil
	}

	fromRate, exists := rates.Rates[fromCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not found in exchange rates", fromCurrency)
	}

	toRate, exists := rates.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not found in exchange rates", toCurrency)
	}

	return amount * (toRate / fromRate), nil
}

func main() {
	apiKey := os.Getenv("EXCHANGE_RATES_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set EXCHANGE_RATES_API_KEY environment variable")
		return
	}

	rates, err := fetchExchangeRates(apiKey)
	if err != nil {
		fmt.Printf("Error fetching exchange rates: %v\n", err)
		return
	}

	if len(os.Args) != 4 {
		fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
		fmt.Println("Example: go run currency_converter.go 100 USD EUR")
		return
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		return
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	convertedAmount, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Error converting currency: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s (as of %s)\n", amount, fromCurrency, convertedAmount, toCurrency, rates.Date)
}