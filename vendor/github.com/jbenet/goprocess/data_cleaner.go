
package main

import (
	"fmt"
	"sort"
)

type DataRecord struct {
	ID   int
	Name string
}

type DataSet []DataRecord

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

func CleanData(records DataSet) DataSet {
	unique := records.RemoveDuplicates()
	unique.SortByID()
	return unique
}

func main() {
	data := DataSet{
		{ID: 5, Name: "ItemE"},
		{ID: 2, Name: "ItemB"},
		{ID: 5, Name: "ItemE"},
		{ID: 1, Name: "ItemA"},
		{ID: 2, Name: "ItemB"},
	}

	cleaned := CleanData(data)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}