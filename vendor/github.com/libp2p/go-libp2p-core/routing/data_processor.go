package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func CleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run data_processor.go <input.csv> <output.csv>")
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := CleanCSV(inputFile, outputFile); err != nil {
		log.Fatal("Error processing CSV:", err)
	}
	log.Println("CSV processing completed successfully")
}