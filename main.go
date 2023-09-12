package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var extraFieldTypes = map[string]string{
	"one_time_fee_name_": "one_time_fee_price_",
	"monthly_fee_name_":  "monthly_fee_price_",
}

type AdditionalCharges struct {
	FieldNumber int
	ChargeName  string
	ChargeValue string
}

type BroadbandData struct {
	CompanyName                  string
	DiscountsAndBundlesURL       string
	AcpEnabled                   string
	CustomerSupportURL           string
	CustomerSupportPhone         string
	NetworkManagementURL         string
	PrivacyPolicyURL             string
	FccID                        string
	DataServiceID                string
	DataServiceName              string
	FixedOrMobile                string
	DataServicePrice             string
	MonthlyPrice                 string
	BillingFrequencyInMonths     string
	IntroductoryRate             bool
	IntroductoryPeriodInMonths   string
	IntroductoryPricePerMonth    string
	ContractDuration             string
	ContractURL                  string
	EarlyTerminationFee          string
	DLSpeedInKbps                string
	CalculatedDLSpeedInMbps      string
	ULSpeedInKbps                string
	LatencyInMs                  string
	CalculatedULSpeedInMbps      string
	DataIncludedInMonthlyPriceGB string
	OverageFee                   string
	OverageDataAmount            string
	ExtraMonthlyFields           []AdditionalCharges
	ExtraOneTimeFields           []AdditionalCharges
}

var csvFileName string
var outputDirectory string

func main() {
	flag.StringVar(&csvFileName, "inputcsv", "bcd.csv", "the name of the csv file to convert")
	flag.StringVar(&outputDirectory, "outputdir", "./generated-labels", "the name of the directory to output the generated files to")
	uldlAreInMbps := flag.Bool("uldlmbps", false, "interpret ul_speed_in_kbps and dl_speed_in_kbps field in the csv file as Mbps not Kbps (no conversions will be done)")
	checkCsv := flag.Bool("checkcsv", false, "check the csv file for errors (basic checks)")
	flag.Parse()

	if *checkCsv {
		err := checkCsvRecords(uldlAreInMbps)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		fmt.Println("error:", err)
		return
	}

	if _, err := os.Stat(csvFileName); os.IsNotExist(err) {
		fmt.Println("error:", err)
		return
	}

	records, err := loadCSV(csvFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var broadbandData []map[string]string

	header := records[0]

	for _, record := range records[1:] {
		data := make(map[string]string)
		for i, value := range record {
			data[header[i]] = value
		}
		broadbandData = append(broadbandData, data)
	}

	var templateData []BroadbandData
	for _, data := range broadbandData {
		templateEntry := BroadbandData{
			CompanyName:                  data["company_name"], //
			DiscountsAndBundlesURL:       data["discounts_and_bundles_url"],
			AcpEnabled:                   data["acp"],
			CustomerSupportURL:           data["customer_support_url"],
			CustomerSupportPhone:         data["customer_support_phone"],
			NetworkManagementURL:         data["network_management_url"],
			PrivacyPolicyURL:             data["privacy_policy_url"],
			FccID:                        data["fcc_id"],
			DataServiceID:                data["data_service_id"],
			DataServiceName:              data["data_service_name"],             //
			FixedOrMobile:                data["fixed_or_mobile"],               //
			DataServicePrice:             data["data_service_price"],            //
			BillingFrequencyInMonths:     data["billing_frequency_in_months"],   //
			IntroductoryPeriodInMonths:   data["introductory_period_in_months"], //
			IntroductoryPricePerMonth:    data["introductory_price_per_month"],  //
			ContractDuration:             data["contract_duration"],             //
			ContractURL:                  data["contract_url"],                  //
			EarlyTerminationFee:          data["early_termination_fee"],
			DLSpeedInKbps:                data["dl_speed_in_kbps"],
			ULSpeedInKbps:                data["ul_speed_in_kbps"],
			LatencyInMs:                  data["latency_in_ms"],
			DataIncludedInMonthlyPriceGB: data["data_included_in_monthly_price"],
			OverageFee:                   data["overage_fee"],
			OverageDataAmount:            data["overage_data_amount"],
		}

		err := calculateUploadDownloadSpeeds(&templateEntry, uldlAreInMbps)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if templateEntry.FixedOrMobile == "" {
			templateEntry.FixedOrMobile = "Fixed"
		}

		err = calculateMonthlyPrice(&templateEntry)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		for fieldName, fieldValue := range data {
			if len(fieldName) > 0 {
				for extraFieldKey, extraFieldType := range extraFieldTypes {
					regex := regexp.MustCompile(extraFieldKey)
					if regex.MatchString(fieldName) {
						splitKey := regexp.MustCompile("_").Split(fieldName, -1)
						indexNumber, err := strconv.Atoi(splitKey[len(splitKey)-1])
						if err == nil {
							if fieldValue != "" {
								indexStr := strconv.Itoa(indexNumber)
								if _, ok := data[extraFieldType+indexStr]; !ok {
									fmt.Println("error: missing associated field for", fieldName)
									continue
								}

								if data[extraFieldType+indexStr] == "" {
									fmt.Println("error: empty value for", fieldName)
									continue
								}

								if len(fieldValue) > 36 {
									fieldValue = fieldValue[:33] + "..."
								}

								e := AdditionalCharges{
									FieldNumber: indexNumber,
									ChargeName:  fieldValue,
									ChargeValue: data[extraFieldType+indexStr],
								}

								if strings.Contains(extraFieldType, "one_time") {
									templateEntry.ExtraOneTimeFields = append(templateEntry.ExtraOneTimeFields, e)
									continue
								}
								if strings.Contains(extraFieldType, "monthly") {
									templateEntry.ExtraMonthlyFields = append(templateEntry.ExtraMonthlyFields, e)
									continue
								}
							}
						} else {
							fmt.Println("error converting index number:", err)
						}
					}
				}
			}
		}
		templateData = append(templateData, templateEntry)
	}

	for _, data := range templateData {
		sort.Slice(data.ExtraMonthlyFields, func(i, j int) bool {
			return strings.ToLower(data.ExtraMonthlyFields[i].ChargeName) < strings.ToLower(data.ExtraMonthlyFields[j].ChargeName)
		})
		sort.Slice(data.ExtraOneTimeFields, func(i, j int) bool {
			return strings.ToLower(data.ExtraOneTimeFields[i].ChargeName) < strings.ToLower(data.ExtraOneTimeFields[j].ChargeName)
		})
	}

	generateLabels(templateData)
}
