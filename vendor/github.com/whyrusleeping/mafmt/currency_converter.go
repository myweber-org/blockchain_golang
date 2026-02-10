
package main

import (
	"fmt"
	"os"
	"strconv"
)

const usdToEurRate = 0.85

func convertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		os.Exit(1)
	}

	amountStr := os.Args[1]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("Invalid amount: %s\n", amountStr)
		os.Exit(1)
	}

	if amount < 0 {
		fmt.Println("Amount cannot be negative")
		os.Exit(1)
	}

	converted := convertUSDToEUR(amount)
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
}