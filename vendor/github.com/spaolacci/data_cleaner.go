package main

import (
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Name string
}

func removeDuplicates(records []Record) []Record {
	seen := make(map[int]bool)
	var result []Record
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func sortRecords(records []Record) []Record {
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})
	return records
}

func cleanData(records []Record) []Record {
	unique := removeDuplicates(records)
	sorted := sortRecords(unique)
	return sorted
}

func main() {
	data := []Record{
		{ID: 3, Name: "Charlie"},
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "Alice"},
		{ID: 4, Name: "David"},
	}

	cleaned := cleanData(data)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", record.ID, record.Name)
	}
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}