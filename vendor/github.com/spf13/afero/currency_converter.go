
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}package main

import (
	"fmt"
	"os"
	"strconv"
)

const exchangeRate = 0.85

func convertUSDToEUR(amount float64) float64 {
	return amount * exchangeRate
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		return
	}

	amountStr := os.Args[1]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("Invalid amount: %s\n", amountStr)
		return
	}

	if amount < 0 {
		fmt.Println("Amount cannot be negative")
		return
	}

	converted := convertUSDToEUR(amount)
	fmt.Printf("%.2f USD = %.2f EUR (Rate: %.2f)\n", amount, converted, exchangeRate)
}