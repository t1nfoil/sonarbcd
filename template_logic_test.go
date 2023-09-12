package main

import (
	"errors"
	"testing"
)

func TestCalculateMonthlyPrice(t *testing.T) {
	testCases := []struct {
		name           string
		templateData   BroadbandData
		expectedResult string
		expectedError  error
	}{
		{
			name: "No Introductory Period",
			templateData: BroadbandData{
				BillingFrequencyInMonths: "12",
				DataServicePrice:         "$100.00",
			},
			expectedResult: "1200.00",
			expectedError:  nil,
		},
		{
			name: "With Introductory Period",
			templateData: BroadbandData{
				BillingFrequencyInMonths:  "12",
				IntroductoryPricePerMonth: "$80.00",
			},
			expectedResult: "960.00",
			expectedError:  nil,
		},
		{
			name: "Invalid Billing Frequency",
			templateData: BroadbandData{
				BillingFrequencyInMonths: "invalid",
				DataServicePrice:         "$100.00",
			},
			expectedResult: "",
			expectedError:  errors.New(`strconv.Atoi: parsing "invalid": invalid syntax`),
		},
		{
			name: "Invalid Price Format",
			templateData: BroadbandData{
				BillingFrequencyInMonths: "12",
				DataServicePrice:         "invalid",
			},
			expectedResult: "",
			expectedError:  errors.New(`strconv.ParseFloat: parsing "invalid": invalid syntax`),
		},
		{
			name: "No Introductory Period, no $",
			templateData: BroadbandData{
				BillingFrequencyInMonths: "12",
				DataServicePrice:         "100.00",
			},
			expectedResult: "1200.00",
			expectedError:  nil,
		},
		{
			name: "With Introductory Period, no $",
			templateData: BroadbandData{
				BillingFrequencyInMonths:  "12",
				IntroductoryPricePerMonth: "80.00",
			},
			expectedResult: "960.00",
			expectedError:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := calculateMonthlyPrice(&testCase.templateData)

			if err != nil {
				if testCase.expectedError == nil {
					t.Errorf("Expected no error, got: %v", err)
				} else if err.Error() != testCase.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", testCase.expectedError, err)
				}
			}

			if testCase.expectedResult != testCase.templateData.MonthlyPrice {
				t.Errorf("Expected result: %s, got: %s", testCase.expectedResult, testCase.templateData.MonthlyPrice)
			}
		})
	}
}

func TestCalculateUploadDownloadSpeeds(t *testing.T) {
	tests := []struct {
		description     string
		templateEntry   BroadbandData
		uldlAreInMbps   bool
		expectedULSpeed string
		expectedDLSpeed string
		expectedError   bool
	}{
		{
			description:     "Test with uldlAreInMbps=true",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10", DLSpeedInKbps: "20"},
			uldlAreInMbps:   true,
			expectedULSpeed: "10.00",
			expectedDLSpeed: "20.00",
			expectedError:   false,
		},
		{
			description:     "Test with uldlAreInMbps=false",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000", DLSpeedInKbps: "20000"},
			uldlAreInMbps:   false,
			expectedULSpeed: "10.00",
			expectedDLSpeed: "20.00",
			expectedError:   false,
		},
		{
			description:     "Test with uldlAreInMbps=true, non-numeric ULSpeedInKbps",
			templateEntry:   BroadbandData{ULSpeedInKbps: "not a number", DLSpeedInKbps: "20000"},
			uldlAreInMbps:   true,
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with uldlAreInMbps=false, non-numeric DLSpeedInKbps",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000", DLSpeedInKbps: "not a number"},
			uldlAreInMbps:   false,
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with uldlAreInMbps=true, empty strings",
			templateEntry:   BroadbandData{ULSpeedInKbps: "", DLSpeedInKbps: ""},
			uldlAreInMbps:   true,
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with uldlAreInMbps=false, empty strings",
			templateEntry:   BroadbandData{ULSpeedInKbps: "", DLSpeedInKbps: ""},
			uldlAreInMbps:   false,
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with uldlAreInMbps=true, valid values at lower boundary",
			templateEntry:   BroadbandData{ULSpeedInKbps: "0", DLSpeedInKbps: "0"},
			uldlAreInMbps:   true,
			expectedULSpeed: "0.00",
			expectedDLSpeed: "0.00",
			expectedError:   false,
		},
		{
			description:     "Test with uldlAreInMbps=false, valid values at lower boundary",
			templateEntry:   BroadbandData{ULSpeedInKbps: "0", DLSpeedInKbps: "0"},
			uldlAreInMbps:   false,
			expectedULSpeed: "0.00",
			expectedDLSpeed: "0.00",
			expectedError:   false,
		},
		{
			description:     "Test with uldlAreInMbps=true, valid values at upper boundary",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000", DLSpeedInKbps: "10000.00"},
			uldlAreInMbps:   true,
			expectedULSpeed: "10000.00",
			expectedDLSpeed: "10000.00",
			expectedError:   false,
		},
		{
			description:     "Test with uldlAreInMbps=false, valid values at upper boundary",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000000", DLSpeedInKbps: "10000000"},
			uldlAreInMbps:   false,
			expectedULSpeed: "10000.00",
			expectedDLSpeed: "10000.00",
			expectedError:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := calculateUploadDownloadSpeeds(&test.templateEntry, &test.uldlAreInMbps)

			if test.expectedError && err == nil {
				t.Errorf("Expected an error but got none")
			}

			if !test.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if test.templateEntry.CalculatedULSpeedInMbps != test.expectedULSpeed {
				t.Errorf("ULSpeed mismatch. Got %s, expected %s", test.templateEntry.CalculatedULSpeedInMbps, test.expectedULSpeed)
			}

			if test.templateEntry.CalculatedDLSpeedInMbps != test.expectedDLSpeed {
				t.Errorf("DLSpeed mismatch. Got %s, expected %s", test.templateEntry.CalculatedDLSpeedInMbps, test.expectedDLSpeed)
			}
		})
	}
}
