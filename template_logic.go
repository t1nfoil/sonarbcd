package main

import (
	"fmt"
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

func calculateUploadDownloadSpeeds(templateEntry *BroadbandData, uldlAreInMbps *bool) error {
	conversionFactor := 0.001
	if *uldlAreInMbps {
		conversionFactor = 1.00

		ulSpeed, err := strconv.ParseFloat(templateEntry.ULSpeedInKbps, 64)
		if err != nil {
			return err
		}

		dlSpeed, err := strconv.ParseFloat(templateEntry.DLSpeedInKbps, 64)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		templateEntry.CalculatedULSpeedInMbps = fmt.Sprintf("%.2f", ulSpeed*conversionFactor)
		templateEntry.CalculatedDLSpeedInMbps = fmt.Sprintf("%.2f", float64(dlSpeed)*conversionFactor)

	} else {
		ulSpeed, err := strconv.ParseInt(templateEntry.ULSpeedInKbps, 10, 64)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		dlSpeed, err := strconv.ParseInt(templateEntry.DLSpeedInKbps, 10, 64)
		if err != nil {
			fmt.Println("error:", err)
			return err
		}

		templateEntry.CalculatedULSpeedInMbps = fmt.Sprintf("%.2f", float64(ulSpeed)*conversionFactor)
		templateEntry.CalculatedDLSpeedInMbps = fmt.Sprintf("%.2f", float64(dlSpeed)*conversionFactor)

	}
	return nil
}
