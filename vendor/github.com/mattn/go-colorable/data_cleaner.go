
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID   string
	Data string
	Hash string
}

func generateHash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		if !seen[record.Hash] {
			seen[record.Hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecord(record Record) bool {
	if strings.TrimSpace(record.ID) == "" {
		return false
	}
	if strings.TrimSpace(record.Data) == "" {
		return false
	}
	expectedHash := generateHash(record.Data)
	return record.Hash == expectedHash
}

func cleanDataset(records []Record) []Record {
	var validRecords []Record
	for _, record := range records {
		if validateRecord(record) {
			validRecords = append(validRecords, record)
		}
	}
	return deduplicateRecords(validRecords)
}

func main() {
	sampleData := []Record{
		{ID: "1", Data: "Sample record one", Hash: generateHash("Sample record one")},
		{ID: "2", Data: "Sample record two", Hash: generateHash("Sample record two")},
		{ID: "1", Data: "Sample record one", Hash: generateHash("Sample record one")},
		{ID: "3", Data: "", Hash: generateHash("")},
		{ID: "", Data: "Invalid record", Hash: generateHash("Invalid record")},
	}

	fmt.Println("Original records:", len(sampleData))
	cleaned := cleanDataset(sampleData)
	fmt.Println("Cleaned records:", len(cleaned))
}