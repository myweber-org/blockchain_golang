package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, line := range input {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	cleaned := cleanData(lines)
	for _, line := range cleaned {
		fmt.Println(line)
	}
}