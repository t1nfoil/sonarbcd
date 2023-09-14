package main

import (
	"errors"
	"strconv"
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

func TestIsSpeedAnInteger(t *testing.T) {
	tests := []struct {
		speed     float64
		isInteger bool
	}{
		{100.00, true},    // Integer, should return true
		{100.5, false},    // Float, should return false
		{0.0, true},       // Zero, should return true
		{-10.5, false},    // Negative float, should return false
		{-20.0, true},     // Negative integer, should return true
		{999999999, true}, // Very large integer, should return true
	}

	for _, test := range tests {
		t.Run(strconv.FormatFloat(test.speed, 'f', -1, 64), func(t *testing.T) {
			result := isSpeedAnInteger(test.speed)
			if result != test.isInteger {
				t.Errorf("isSpeedAnInteger(%f) = %v; want %v", test.speed, result, test.isInteger)
			}
		})
	}
}

func TestCalculateUploadDownloadSpeeds(t *testing.T) {
	tests := []struct {
		description     string
		templateEntry   BroadbandData
		expectedULSpeed string
		expectedDLSpeed string
		expectedError   bool
	}{
		{
			description:     "Test ULSpeedInKbps and DLSpeedInKbps conversion exact decimal to integer",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10.0", DLSpeedInKbps: "20.0"},
			expectedULSpeed: "10",
			expectedDLSpeed: "20",
			expectedError:   false,
		},
		{
			description:     "Test with integer ULSpeedInKbps and DLSpeedInKbps speeds (Kbps, no conversion)",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000", DLSpeedInKbps: "20000"},
			expectedULSpeed: "10",
			expectedDLSpeed: "20",
			expectedError:   false,
		},
		{
			description:     "Test with non-numeric ULSpeedInKbps",
			templateEntry:   BroadbandData{ULSpeedInKbps: "not a number", DLSpeedInKbps: "20000"},
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with non-numeric DLSpeedInKbps",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000", DLSpeedInKbps: "not a number"},
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with empty strings",
			templateEntry:   BroadbandData{ULSpeedInKbps: "", DLSpeedInKbps: ""},
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with uldlAreInMbps=false, empty strings",
			templateEntry:   BroadbandData{ULSpeedInKbps: "", DLSpeedInKbps: ""},
			expectedULSpeed: "",
			expectedDLSpeed: "",
			expectedError:   true,
		},
		{
			description:     "Test with valid values at lower boundary (integer)",
			templateEntry:   BroadbandData{ULSpeedInKbps: "0", DLSpeedInKbps: "0"},
			expectedULSpeed: "0",
			expectedDLSpeed: "0",
			expectedError:   false,
		},
		{
			description:     "Test with valid values at lower boundary (decimal)",
			templateEntry:   BroadbandData{ULSpeedInKbps: "0.00", DLSpeedInKbps: "0.00"},
			expectedULSpeed: "0",
			expectedDLSpeed: "0",
			expectedError:   false,
		},
		{
			description:     "Test with valid values at upper boundary (decimal)",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000.00", DLSpeedInKbps: "10000.00"},
			expectedULSpeed: "10000",
			expectedDLSpeed: "10000",
			expectedError:   false,
		},
		{
			description:     "Test with valid values at upper boundary (integer)",
			templateEntry:   BroadbandData{ULSpeedInKbps: "10000000", DLSpeedInKbps: "10000000"},
			expectedULSpeed: "10000",
			expectedDLSpeed: "10000",
			expectedError:   false,
		},
		{
			description:     "Test integer to decimal conversion test",
			templateEntry:   BroadbandData{ULSpeedInKbps: "1500", DLSpeedInKbps: "1500"},
			expectedULSpeed: "1.5",
			expectedDLSpeed: "1.5",
			expectedError:   false,
		},
		{
			description:     "Test decimal precision conversion",
			templateEntry:   BroadbandData{ULSpeedInKbps: "1.500", DLSpeedInKbps: "1.500"},
			expectedULSpeed: "1.5",
			expectedDLSpeed: "1.5",
			expectedError:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := calculateUploadDownloadSpeeds(&test.templateEntry)

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
