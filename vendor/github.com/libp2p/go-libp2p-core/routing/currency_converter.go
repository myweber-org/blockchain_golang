package main

import (
	"fmt"
	"os"
	"strconv"
)

const exchangeRate = 0.85

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		return
	}

	usdAmount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		return
	}

	if usdAmount < 0 {
		fmt.Println("Amount cannot be negative")
		return
	}

	eurAmount := usdAmount * exchangeRate
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}