
package main

import (
	"fmt"
	"sort"
)

type DataRecord struct {
	ID   int
	Name string
}

func CleanData(records []DataRecord) []DataRecord {
	seen := make(map[int]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ID < unique[j].ID
	})

	return unique
}

func main() {
	records := []DataRecord{
		{ID: 3, Name: "Charlie"},
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "AliceDuplicate"},
		{ID: 4, Name: "David"},
	}

	cleaned := CleanData(records)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}