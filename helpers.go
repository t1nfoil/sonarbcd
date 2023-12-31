package main

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
)

func convertPriceToCents(price string) (int, error) {
	price = strings.TrimPrefix(price, "$")
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, err
	}

	priceCents := int(math.Round(priceFloat * 100))
	return priceCents, nil
}

type jsonError struct {
	IsError string `json:"isError"`
	Message string `json:"message"`
	Row     string `json:"row"`
}

func convertErrorToJSON(row string, messages ...string) error {
	var j jsonError
	j.IsError = "true"
	if row == "NA" {
		row = ""
	}
	j.Row = row
	j.Message = strings.Join(messages, " ")

	json, _ := json.Marshal(j)
	return errors.New(string(json))
}
