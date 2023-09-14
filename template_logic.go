package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func calculateMonthlyPrice(templateEntry *BroadbandData) error {
	templateEntry.DataServicePrice = strings.TrimPrefix(templateEntry.DataServicePrice, "$")

	billingFrequencyInMonths, err := strconv.Atoi(templateEntry.BillingFrequencyInMonths)
	if err != nil {
		return err
	}

	var priceInCents int
	if templateEntry.IntroductoryPeriodInMonths == "" && templateEntry.IntroductoryPricePerMonth == "" {
		priceInCents, err = convertPriceToCents(templateEntry.DataServicePrice)
		if err != nil {
			return err
		}
	} else {
		templateEntry.IntroductoryRate = true
		priceInCents, err = convertPriceToCents(templateEntry.IntroductoryPricePerMonth)
		if err != nil {
			return err
		}
	}

	templateEntry.MonthlyPrice = fmt.Sprintf("%.2f", float64(billingFrequencyInMonths*priceInCents)/100)
	return nil
}

func isSpeedAnInteger(speed float64) bool {
	str := strconv.FormatFloat(speed, 'f', -1, 64)
	isInteger := str == strconv.Itoa(int(speed))
	if isInteger {
		return true
	}
	return false
}

func calculateUploadDownloadSpeeds(templateEntry *BroadbandData) error {

	if match, _ := regexp.MatchString(`[^0-9\.]`, templateEntry.ULSpeedInKbps); match {
		templateEntry.ULSpeedInKbps = ""
		templateEntry.DLSpeedInKbps = ""
		return fmt.Errorf("ULSpeedInKbps contains invalid characters")
	}

	if match, _ := regexp.MatchString(`[^0-9\.]`, templateEntry.DLSpeedInKbps); match {
		templateEntry.ULSpeedInKbps = ""
		templateEntry.DLSpeedInKbps = ""
		return fmt.Errorf("ULSpeedInKbps contains invalid characters")
	}

	if strings.Contains(templateEntry.ULSpeedInKbps, ".") {
		conversionFactor := 1.00
		ulSpeed, err := strconv.ParseFloat(templateEntry.ULSpeedInKbps, 64)
		if err != nil {
			return err
		}

		conversionFormat := "%.1f"
		if isSpeedAnInteger(ulSpeed) {
			conversionFormat = "%.0f"
		}
		templateEntry.CalculatedULSpeedInMbps = fmt.Sprintf(conversionFormat, ulSpeed*conversionFactor)
	} else {
		conversionFactor := 0.001
		ulSpeed, err := strconv.ParseInt(templateEntry.ULSpeedInKbps, 10, 64)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		ulSpeedMbps := float64(ulSpeed) * conversionFactor
		if ulSpeedMbps == float64(int(ulSpeedMbps)) {
			templateEntry.CalculatedULSpeedInMbps = fmt.Sprintf("%.0f", ulSpeedMbps)
		} else {
			templateEntry.CalculatedULSpeedInMbps = fmt.Sprintf("%.1f", ulSpeedMbps)
		}
	}

	if strings.Contains(templateEntry.DLSpeedInKbps, ".") {
		conversionFactor := 1.00
		dlSpeed, err := strconv.ParseFloat(templateEntry.DLSpeedInKbps, 64)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		conversionFormat := "%.1f"
		if isSpeedAnInteger(dlSpeed) {
			conversionFormat = "%.0f"
		}
		templateEntry.CalculatedDLSpeedInMbps = fmt.Sprintf(conversionFormat, float64(dlSpeed)*conversionFactor)

	} else {
		conversionFactor := 0.001
		dlSpeed, err := strconv.ParseInt(templateEntry.DLSpeedInKbps, 10, 64)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		dlSpeedMbps := float64(dlSpeed) * conversionFactor
		if dlSpeedMbps == float64(int(dlSpeedMbps)) {
			templateEntry.CalculatedDLSpeedInMbps = fmt.Sprintf("%.0f", dlSpeedMbps)
		} else {
			templateEntry.CalculatedDLSpeedInMbps = fmt.Sprintf("%.1f", dlSpeedMbps)
		}
	}

	return nil
}
