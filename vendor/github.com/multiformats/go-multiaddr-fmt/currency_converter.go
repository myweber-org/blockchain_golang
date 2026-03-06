
package main

import (
	"fmt"
)

func main() {
	const usdToEurRate = 0.85
	var usdAmount float64

	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	eurAmount := usdAmount * usdToEurRate
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}