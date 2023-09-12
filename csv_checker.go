package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func checkCsvRecords(uldlAreInMbps *bool) error {
	records, err := loadCSV(csvFileName)
	if err != nil {
		return err
	}
	header := records[0]

	for recordNumber, record := range records[1:] {
		data := make(map[string]string)
		for i, value := range record {
			data[header[i]] = value
		}
		data["csvrow"] = strconv.Itoa(recordNumber + 2)
		validateIntroductoryFields(data)
		validateDataServicePrice(data)
		validateSpeeds(data, *uldlAreInMbps)
	}

	return nil
}

func validateIntroductoryFields(data map[string]string) {
	introductoryPeriod := data["introductory_period_in_months"]
	introductoryPrice := data["introductory_price_per_month"]

	if introductoryPeriod != "" || introductoryPrice != "" {
		// ensure they both are present
		if introductoryPeriod == "" || introductoryPrice == "" {
			fmt.Println("CSV: Introductory period and price must both be present if either are set, row ", data["csvrow"])
			return
		}

		if _, err := strconv.Atoi(fmt.Sprintf("%v", introductoryPeriod)); err != nil {
			fmt.Println("CSV: Introductory period must be a valid integer, row ", data["csvrow"])
			return
		}

		price := fmt.Sprintf("%v", introductoryPrice)
		if match, _ := regexp.MatchString(`^\$?\d{1,3}(\.\d{1,2})?$`, price); !match || len(price) > 8 {
			fmt.Println("CSV: Introductory price format should be [$]###.##, row ", data["csvrow"])
			return
		}

		priceValue, err := strconv.ParseFloat(strings.TrimLeft(price, "$"), 64)
		if err != nil {
			fmt.Println("CSV: Introductory price could not be converted to float64, row ", data["csvrow"])
			return
		} else if priceValue != float64(int64(priceValue*100))/100 {
			fmt.Println("CSV: Introductory price must have 2 decimal precision, row ", data["csvrow"])
			return
		}
	}
}

func validateDataServicePrice(data map[string]string) {
	dataServicePrice, exists := data["data_service_price"]

	if exists {
		price := fmt.Sprintf("%v", dataServicePrice)
		if match, _ := regexp.MatchString(`^\$?\d{1,3}(\.\d{3})*(\.\d{1,3})?$`, price); !match || len(price) > 8 {
			fmt.Println("CSV: Data service price format should be [$]###.###, row ", data["csvrow"])
			return
		}

		priceValue, err := strconv.ParseFloat(strings.TrimLeft(price, "$"), 64)
		if err != nil {
			fmt.Println("CSV: Data service price could not be converted to float64, row ", data["csvrow"])
		} else if priceValue != float64(int64(priceValue*1000))/1000 {
			fmt.Println("CSV: Data service price must have 3 decimal precision, row ", data["csvrow"])
			return
		}
	}
}

func validateSpeeds(data map[string]string, uldlAreInMbps bool) {
	dlSpeed, dlExists := data["dl_speed_in_kbps"]
	ulSpeed, ulExists := data["ul_speed_in_kbps"]

	if dlExists && ulExists && !uldlAreInMbps {
		dlSpeedValue, dlErr := strconv.ParseInt(fmt.Sprintf("%v", dlSpeed), 10, 64)
		ulSpeedValue, ulErr := strconv.ParseInt(fmt.Sprintf("%v", ulSpeed), 10, 64)

		if dlErr != nil || ulErr != nil {
			fmt.Println("Error: Speed values must be valid integers, row ", data["csvrow"])
			return
		}

		if dlSpeedValue < 0 || dlSpeedValue > 10000000 || ulSpeedValue < 0 || ulSpeedValue > 10000000 {
			fmt.Println("Error: Speed values must be between 0 and 10000000, row ", data["csvrow"])
			return
		}
	} else {
		dlSpeedValue, dlErr := strconv.ParseFloat(fmt.Sprintf("%v", dlSpeed), 64)
		ulSpeedValue, ulErr := strconv.ParseFloat(fmt.Sprintf("%v", ulSpeed), 64)
		if dlErr != nil || ulErr != nil {
			fmt.Println("Error: Speed values must be valid integers or decimal values, when using -uldlmbps flag, row ", data["csvrow"])
			return
		}

		if dlSpeedValue < 0 || dlSpeedValue > 10000 || ulSpeedValue < 0 || ulSpeedValue > 10000 {
			fmt.Println("Error: Speed values must be between 0.00 and 10000, row ", data["csvrow"])
			return
		}

		if dlSpeedValue < 0 || dlSpeedValue > 10000 || ulSpeedValue < 0 || ulSpeedValue > 10000 {
			fmt.Println("Error: Speed values must be between 0 and 10000, row ", data["csvrow"])
			return
		}
	}
}
