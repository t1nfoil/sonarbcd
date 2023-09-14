package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func checkCsvRecords() error {
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
		err := validateIntroductoryFields(data)
		if err != nil {
			return err
		}

		err = validateDataServicePrice(data)
		if err != nil {
			return err
		}

		err = validateSpeeds(data)
		if err != nil {
			return err
		}

	}

	return nil
}

func validateIntroductoryFields(data map[string]string) error {
	introductoryPeriod := data["introductory_period_in_months"]
	introductoryPrice := data["introductory_price_per_month"]

	if introductoryPeriod != "" || introductoryPrice != "" {
		if introductoryPeriod == "" || introductoryPrice == "" {
			return errors.New("CSV: Introductory period and price must both be present if either are set, row " + data["csvrow"])

		}

		if _, err := strconv.Atoi(fmt.Sprintf("%v", introductoryPeriod)); err != nil {
			return errors.New("CSV: Introductory period must be a valid integer, row " + data["csvrow"])
		}

		price := fmt.Sprintf("%v", introductoryPrice)
		if match, _ := regexp.MatchString(`^\$?\d{1,3}(\.\d{1,2})?$`, price); !match || len(price) > 8 {
			return errors.New("CSV: Introductory price format should be [$]###.##, row " + data["csvrow"])
		}

		priceValue, err := strconv.ParseFloat(strings.TrimLeft(price, "$"), 64)
		if err != nil {
			return errors.New("CSV: Introductory price could not be converted to float64, row " + data["csvrow"])
		} else if priceValue != float64(int64(priceValue*100))/100 {
			return errors.New("CSV: Introductory price must have 2 decimal precision, row " + data["csvrow"])
		}
	}
	return nil
}

func validateDataServicePrice(data map[string]string) error {
	dataServicePrice, exists := data["data_service_price"]

	if exists {
		price := fmt.Sprintf("%v", dataServicePrice)
		if match, _ := regexp.MatchString(`^\$?\d{1,3}(\.\d{3})*(\.\d{1,3})?$`, price); !match || len(price) > 8 {
			return errors.New("CSV: Data service price format should be [$]###.###, row " + data["csvrow"])
		}

		priceValue, err := strconv.ParseFloat(strings.TrimLeft(price, "$"), 64)
		if err != nil {
			return errors.New("CSV: Data service price could not be converted to float64, row " + data["csvrow"])

		} else if priceValue != float64(int64(priceValue*1000))/1000 {
			return errors.New("CSV: Data service price must have 3 decimal precision, row " + data["csvrow"])
		}
	}
	return nil
}

func validateSpeeds(data map[string]string) error {
	dlSpeed, dlExists := data["dl_speed_in_kbps"]
	ulSpeed, ulExists := data["ul_speed_in_kbps"]

	if dlExists && strings.Contains(dlSpeed, ".") {
		dlSpeedValue, dlErr := strconv.ParseFloat(fmt.Sprintf("%v", dlSpeed), 64)
		if dlErr != nil {
			return errors.New("error: dl_speed_in_kbps values must be a valid decimal value to be interpreted as Mpbs, row" + data["csvrow"])
		}

		if dlSpeedValue < 0 || dlSpeedValue > 10000 {
			return errors.New("error: dl_speed_in_kbps values must be between 0.00 and 10000.00, row " + data["csvrow"])
		}

	} else {
		dlSpeedValue, dlErr := strconv.ParseInt(fmt.Sprintf("%v", dlSpeed), 10, 64)

		if dlErr != nil {
			return errors.New("error: dl_speed_in_kbps values must be a valid integer (Kbps), row " + data["csvrow"])
		}
		if dlSpeedValue < 0 || dlSpeedValue > 10000000 {
			return errors.New("error: dl_speed_in_kbps values must be between 0 and 10000000, row " + data["csvrow"])
		}

	}

	if ulExists && strings.Contains(ulSpeed, ".") {
		ulSpeedValue, ulErr := strconv.ParseFloat(fmt.Sprintf("%v", ulSpeed), 64)
		if ulErr != nil {
			return errors.New("error: ul_speed_in_kbps values must be a valid decimal value to be interpreted as Mpbs, row " + data["csvrow"])
		}

		if ulSpeedValue < 0 || ulSpeedValue > 10000 {
			return errors.New("error: ul_speed_in_kbps values must be between 0.00 and 10000.00, row " + data["csvrow"])
		}

	} else {
		ulSpeedValue, ulErr := strconv.ParseInt(fmt.Sprintf("%v", ulSpeed), 10, 64)

		if ulErr != nil {
			return errors.New("error: ul_speed_in_kbps values must be valid integer (Kbps), row " + data["csvrow"])
		}

		if ulSpeedValue < 0 || ulSpeedValue > 10000000 {
			return errors.New("error: ul_speed_in_kbps values must be between 0 and 10000000, row " + data["csvrow"])
		}
	}
	return nil
}
