package models

import (
	"fmt"
	"strconv"
	"strings"
)

const Currency = '$'

// Dollars returns a string that represents the value of the currency given by "amount"
func Dollars(amount int) string {
	sign := ""
	negative := amount < 0
	if negative {
		amount *= -1
		sign = "-"
	}
	dollars := amount / 100
	cents := amount % 100
	return fmt.Sprintf("%s%c%d.%02d", sign, Currency, dollars, cents)
}

// Cents takes a currency string formatted as [$]X.XX and returns the number of
// cents that it represents
func Cents(currency string) (int, error) {
	currency = strings.Replace(currency, ".", "", 1)
	currency = strings.Replace(currency, string(Currency), "", 1)
	currency = strings.Replace(currency, ",", "", 1)
	amount, err := strconv.Atoi(currency)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
