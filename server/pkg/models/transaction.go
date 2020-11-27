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

// Transaction represents a single transaction in a person's budget
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
				"(%s TEXT NOT NULL, %s INTEGER NOT NULL, %s INTEGER NOT NULL, %s TEXT NOT NULL)",
			TransactionTableName,
			TransactionEntityCol,
			TransactionAmountCol,
			TransactionDateCol,
			TransactionNoteCol,
		),
	)
}

// InsertTransaction inserts a transaction into the transactions table
func InsertTransaction(db *sql.DB, tx Transaction) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s VALUES (?, ?, ?, ?)",
			TransactionTableName,
		),
		tx.Entity,
		tx.Amount,
		tx.Date,
		tx.Note,
	)
}

// RemoveTransaction removes a transaction from the transactions table
func RemoveTransaction(db *sql.DB, tx Transaction) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE %s=? AND %s=? AND %s=? AND %s=?",
			TransactionTableName,
			TransactionEntityCol,
			TransactionAmountCol,
			TransactionDateCol,
			TransactionNoteCol,
		),
		tx.Entity,
		tx.Amount,
		tx.Date,
		tx.Note,
	)
}

// UpdateTransaction updates a transaction from the transactions table
func UpdateTransaction(db *sql.DB, old, new Transaction) (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET %s=?, %s=? ,%s=?, %s=? WHERE %s=? AND %s=? AND %s=? AND %s=?",
			TransactionTableName,
			TransactionEntityCol,
			TransactionAmountCol,
			TransactionDateCol,
			TransactionNoteCol,
			TransactionEntityCol,
			TransactionAmountCol,
			TransactionDateCol,
			TransactionNoteCol,
		),
		new.Entity,
		new.Amount,
		new.Date,
		new.Note,
		old.Entity,
		old.Amount,
		old.Date,
		old.Note,
	)
}
