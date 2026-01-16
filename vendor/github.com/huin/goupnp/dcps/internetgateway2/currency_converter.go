
package main

import (
    "fmt"
    "os"
    "strconv"
)

const usdToEurRate = 0.85

func main() {
    if len(os.Args) != 2 {
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

    eurAmount := amount * usdToEurRate
    fmt.Printf("%.2f USD = %.2f EUR (Rate: %.2f)\n", amount, eurAmount, usdToEurRate)
}