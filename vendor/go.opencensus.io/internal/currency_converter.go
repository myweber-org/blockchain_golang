
package main

import (
	"fmt"
)

func convertUSDToEUR(amount float64) float64 {
	const exchangeRate = 0.85
	return amount * exchangeRate
}

func main() {
	var usdAmount float64
	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	if usdAmount < 0 {
		fmt.Println("Amount cannot be negative")
		return
	}

	eurAmount := convertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}