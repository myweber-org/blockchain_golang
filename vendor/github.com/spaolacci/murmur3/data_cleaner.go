
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
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

func (dc *DataCleaner) ProcessRecords(records []string) []string {
	cleaned := dc.RemoveDuplicates(records)
	return cleaned
}

func main() {
	cleaner := &DataCleaner{}
	data := []string{"apple", " banana", "apple", "cherry ", "", "banana", "  "}
	result := cleaner.ProcessRecords(data)
	fmt.Println("Cleaned data:", result)
}
package main

import (
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Name string
}

type DataSet []Record

func (d DataSet) RemoveDuplicates() DataSet {
	seen := make(map[int]bool)
	result := DataSet{}
	for _, record := range d {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func (d DataSet) SortByID() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
}

func CleanData(data DataSet) DataSet {
	uniqueData := data.RemoveDuplicates()
	uniqueData.SortByID()
	return uniqueData
}

func main() {
	sampleData := DataSet{
		{3, "Charlie"},
		{1, "Alice"},
		{2, "Bob"},
		{1, "Alice"},
		{4, "David"},
		{2, "Bob"},
	}

	fmt.Println("Original data:", sampleData)
	cleanedData := CleanData(sampleData)
	fmt.Println("Cleaned data:", cleanedData)
}