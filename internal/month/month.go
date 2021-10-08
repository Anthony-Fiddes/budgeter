// Package month provides easy ways to manipulate time.Time's based on the
// context of months.
package month

import "time"

// Start returns the inputted time, but at the beginning of the first day of the
// its month.
func Start(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// End returns the inputted time, but at the end of the last day of its month.
func End(t time.Time) time.Time {
	result := Start(t)
	result = result.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return result
}

func Add(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}
