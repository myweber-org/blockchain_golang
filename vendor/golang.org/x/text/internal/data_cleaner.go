
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    
    return result
}

func main() {
    data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := RemoveDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) TrimWhitespace() []string {
	result := make([]string, len(dc.Data))
	for i, item := range dc.Data {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	rawData := []string{"  apple  ", "banana", "  apple  ", " cherry", "banana "}
	cleaner := NewDataCleaner(rawData)

	trimmed := cleaner.TrimWhitespace()
	cleaner.Data = trimmed
	unique := cleaner.RemoveDuplicates()

	fmt.Println("Original:", rawData)
	fmt.Println("Cleaned:", unique)
}package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	items []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{items: data}
}

func (dc *DataCleaner) RemoveDuplicates() {
	seen := make(map[string]bool)
	var result []string
	for _, item := range dc.items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	dc.items = result
}

func (dc *DataCleaner) TrimWhitespace() {
	for i, item := range dc.items {
		dc.items[i] = strings.TrimSpace(item)
	}
}

func (dc *DataCleaner) GetData() []string {
	return dc.items
}

func main() {
	rawData := []string{"apple ", " banana", "apple", " cherry ", "banana"}
	cleaner := NewDataCleaner(rawData)
	
	cleaner.TrimWhitespace()
	cleaner.RemoveDuplicates()
	
	fmt.Println("Cleaned data:", cleaner.GetData())
}