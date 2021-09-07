// transaction provides a model for a single user transaction. It also provides
// a simple implementation of a sqlite table for storing and querying transactions.
package transaction

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TableName  = "transactions"
	IDCol      = "ID"
	EntityCol  = "Entity"
	AmountCol  = "Amount"
	DateCol    = "Date"
	NoteCol    = "Note"
	DateLayout = "1/2/2006"
	// TODO: this should probably be configurable, but I currently only use US dollars
	Currency  = "$"
	Point     = "."
	Thousands = ","
)

// TODO: add a String() function
// Transaction represents a single transaction in a person's budget
type Transaction struct {
	ID int
	// Entity is the person or company the transaction was made with.
	Entity string
	// Amount is the cost of the transaction in cents
	Amount int
	// Date is the Unix Time the transaction occurred.
	Date int64
	// Note is any note the user wants to add about the transaction.
	Note string
}

// DateString returns the Transaction's date in M/D/YYYY format.
func (t Transaction) DateString() string {
	d := time.Unix(t.Date, 0).UTC()
	date := d.Format(DateLayout)
	return date
}

// AmountString returns a string that represents the amount of the currency that
//the Transaction was worth.
func (t Transaction) AmountString() string {
	return Dollars(t.Amount)
}

// Date converts a string of format M/D/YYYY and converts it to the appropriate
// Unix time. This function is useful for working with the "Transaction" type.
func Date(date string) (int64, error) {
	result, err := time.Parse(DateLayout, date)
	if err != nil {
		return 0, fmt.Errorf("transaction: date \"%s\" must be provided in M/D/YYYY format", date)
	}
	return result.Unix(), nil
}

// Dollars returns a string that represents the value of the given number of
// "cents".
func Dollars(cents int) string {
	sign := ""
	negative := cents < 0
	if negative {
		cents *= -1
		sign = "-"
	}
	dollars := cents / 100
	pennies := cents % 100
	// TODO: determine whether or not I want this to add in thousands separators
	return fmt.Sprintf("%s%s%d.%02d", sign, Currency, dollars, pennies)
}

// Cents takes a currency string formatted as [$]X.XX and returns the number of
// cents that it represents
func Cents(currency string) (int, error) {
	currencyErr := func() error {
		errCurrencyFmt := fmt.Sprintf(
			"transaction: currency \"%%s\" must be provided in [%s]X%sXX format (\"%s\" is allowed)",
			Currency, Point, Thousands,
		)
		return fmt.Errorf(errCurrencyFmt, currency)
	}
	currency = strings.Replace(currency, Currency, "", 1)
	currency = strings.Replace(currency, Thousands, "", 1)
	c := strings.Split(currency, ".")
	dollars, err := strconv.Atoi(c[0])
	if err != nil {
		return 0, currencyErr()
	}
	pennies := 0
	if len(c) == 2 {
		pennies, err = strconv.Atoi(c[1])
		if err != nil {
			return 0, currencyErr()
		}
		if dollars < 0 {
			pennies *= -1
		}
	} else if len(c) > 2 {
		return 0, currencyErr()
	}
	return dollars*100 + pennies, nil
}
