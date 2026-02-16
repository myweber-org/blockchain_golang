package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func fetchExchangeRates(baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", baseCurrency)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rateResponse ExchangeRateResponse
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		return nil, err
	}

	return rateResponse.Rates, nil
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	rates, err := fetchExchangeRates(fromCurrency)
	if err != nil {
		return 0, err
	}

	rate, exists := rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("exchange rate for %s not found", toCurrency)
	}

	return amount * rate, nil
}

func main() {
	if len(os.Args) != 4 {
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

	result, err := convertCurrency(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}