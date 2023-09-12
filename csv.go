package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func loadCSV(csvFile string) ([][]string, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return records, nil
}
