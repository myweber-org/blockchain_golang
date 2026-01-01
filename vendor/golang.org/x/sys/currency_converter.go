
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	var usdAmount float64
	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}