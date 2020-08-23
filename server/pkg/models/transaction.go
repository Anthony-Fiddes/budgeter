package models

type Transaction struct {
	Entity string
	// In 1/100's of a cent
	Amount int
	// Unix Time
	Date int
	Note string
}
