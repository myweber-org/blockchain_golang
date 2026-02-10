package main

import (
	"fmt"
)

const usdToEurRate = 0.92

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 150.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}