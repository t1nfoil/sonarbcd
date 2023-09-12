package main

import (
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
