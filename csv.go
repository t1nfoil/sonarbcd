package main

import (
	"encoding/csv"
	"os"
)

func loadCSV(csvFile string) ([][]string, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
