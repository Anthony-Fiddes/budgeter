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

// Cent represents 1/100th of a Dollar
type Cent int

// String returns a string that represents the value of the given number of
// "cents".
func (c Cent) String() string {
	sign := ""
	negative := c < 0
	if negative {
		c *= -1
		sign = "-"
	}
	dollars := c / 100
	pennies := c % 100
	// TODO: determine whether or not I want this to add in thousands separators
	return fmt.Sprintf("%s%s%d.%02d", sign, Currency, dollars, pennies)
}

// TODO: add a String() function
// Transaction represents a single transaction in a person's budget
type Transaction struct {
	ID int
	// Entity is the person or company the transaction was made with.
	Entity string
	// Amount is the cost of the transaction in cents
	Amount Cent
	// Date is the Unix Time the transaction occurred in seconds.
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

// Unix converts a string of format M/D/YYYY and converts it to the appropriate
// Unix time in seconds. This function is useful for working with the "Transaction" type.
func Unix(date string) (int64, error) {
	result, err := time.Parse(DateLayout, date)
	if err != nil {
		return 0, fmt.Errorf("transaction: date \"%s\" must be provided in M/D/YYYY format", date)
	}
	return result.Unix(), nil
}

// GetCents takes a currency string formatted as [$]X.XX and returns the number of
// cents that it represents
func GetCents(currency string) (Cent, error) {
	currencyErr := func() error {
		errCurrencyFmt := fmt.Sprintf(
			"transaction: currency \"%%s\" must be provided in [%s]X%sXX format (\"%s\" is allowed)",
			Currency, Point, Thousands,
		)
		return fmt.Errorf(errCurrencyFmt, currency)
	}
	currency = strings.Replace(currency, Currency, "", 1)
	currency = strings.Replace(currency, Thousands, "", 1)
	currency = strings.TrimSpace(currency)
	negative := false
	if strings.ContainsRune(currency, '-') {
		negative = true
		currency = strings.Replace(currency, "-", "", 1)
	}
	currStr := strings.Split(currency, ".")
	if currStr[0] == "" {
		currStr[0] = "0"
	}
	dollars, err := strconv.Atoi(currStr[0])
	if err != nil {
		return 0, currencyErr()
	}
	if negative {
		dollars *= -1
	}
	pennies := 0
	if len(currStr) == 2 {
		penStr := currStr[1]
		pennies, err = strconv.Atoi(penStr)
		if len(penStr) == 1 {
			pennies *= 10
		}
		if err != nil {
			return 0, currencyErr()
		}
		if negative {
			pennies *= -1
		}
	} else if len(currStr) > 2 {
		return 0, currencyErr()
	}
	return Cent(dollars*100 + pennies), nil
}
