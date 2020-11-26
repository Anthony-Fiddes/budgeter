package models

import (
	"database/sql"
	"fmt"
)

const (
	TransactionTableName = "transactions"
	TransactionEntityCol = "Entity"
	TransactionAmountCol = "Amount"
	TransactionDateCol   = "Date"
	TransactionNoteCol   = "Note"
)

type Transaction struct {
	Entity string
	// In 1/100's of a cent
	Amount int
	// Unix Time
	Date int
	Note string
}

// CreateTransactionTable creates the transaction table if it does not already exist.
func CreateTransactionTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s "+
				"(id INTEGER PRIMARY KEY, %s TEXT, %s INTEGER, %s INTEGER, %s TEXT)",
			TransactionTableName,
			TransactionEntityCol,
			TransactionAmountCol,
			TransactionDateCol,
			TransactionNoteCol,
		),
	)
}
