package models

import (
	"errors"
	"time"
)

// Date converts a string of format M/D/YYYY and converts it to the appropriate
// Unix time. This function is useful for working with the "Transaction" type.
func Date(date string) (int64, error) {
	result, err := time.Parse(DateLayout, date)
	if err != nil {
		return 0, errors.New("date must be provided in M/D/YYYY format")
	}
	return result.Unix(), nil
}
