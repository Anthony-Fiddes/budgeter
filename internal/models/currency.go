package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// TODO: this should probably be configurable, but I currently only use US dollars
	Currency  = "$"
	Point     = "."
	Thousands = ","
)

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
	// TODO: determine whether or not I want this to add in thousands separators
	return fmt.Sprintf("%s%s%d.%02d", sign, Currency, dollars, cents)
}

// Cents takes a currency string formatted as [$]X.XX and returns the number of
// cents that it represents
func Cents(currency string) (int, error) {
	currencyErr := func() error {
		errCurrencyFmt := fmt.Sprintf("currency \"%%s\" must be provided in [%s]X%sXX format (\"%s\" is allowed)", Currency, Point, Thousands)
		return fmt.Errorf(errCurrencyFmt, currency)
	}
	currency = strings.Replace(currency, Currency, "", 1)
	currency = strings.Replace(currency, Thousands, "", 1)
	c := strings.Split(currency, ".")
	dollars, err := strconv.Atoi(c[0])
	if err != nil {
		return 0, currencyErr()
	}
	cents := 0
	if len(c) == 2 {
		cents, err = strconv.Atoi(c[1])
		if err != nil {
			return 0, currencyErr()
		}
		if dollars < 0 {
			cents *= -1
		}
	} else if len(c) > 2 {
		return 0, currencyErr()
	}
	return dollars*100 + cents, nil
}
