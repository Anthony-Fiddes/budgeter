package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	TransactionTableName = "transactions"
	TransactionEntityCol = "Entity"
	TransactionAmountCol = "Amount"
	TransactionDateCol   = "Date"
	TransactionNoteCol   = "Note"
	DateLayout           = "1/2/2006"
)

// Transaction represents a single transaction in a person's budget
type Transaction struct {
	Entity string
	// In 1/100's of a cent
	Amount int
	// Unix Time
	Date int64
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

// DB wraps *sql.DB to add methods for putting transactions in the transactions table
type DB struct {
	*sql.DB
}

// CreateTransactionTable creates the transactions table if it does not already exist.
func (db *DB) CreateTransactionTable() (sql.Result, error) {
	return db.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s "+
				"(%s TEXT NOT NULL, %s INTEGER NOT NULL, %s INTEGER NOT NULL, %s TEXT NOT NULL, "+
				"UNIQUE(%s,%s,%s,%s))",
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
	)
}

// GetTransactions gets transactions from the transactions table. The
// most recent transactions will be returned first. GetTransactions will
// return, at most, "limit" results.
func (db *DB) GetTransactions(limit int) ([]Transaction, error) {
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT * FROM %s ORDER BY %s DESC LIMIT ?",
			TransactionTableName,
			TransactionDateCol,
		),
		limit,
	)
	if err != nil {
		return nil, err
	}
	result := make([]Transaction, 0, limit)
	for rows.Next() {
		tx := Transaction{}
		err := rows.Scan(&tx.Entity, &tx.Amount, &tx.Date, &tx.Note)
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	return result, nil
}

// Search returns the most recent transactions that include the given "query".
// It returns, at most, "limit" transactions.
func (db *DB) Search(query string, limit int) ([]Transaction, error) {
	query = "%" + query + "%"
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT * FROM %s WHERE %s LIKE ? OR %s LIKE ? ORDER BY %s DESC LIMIT ?",
			TransactionTableName,
			TransactionEntityCol,
			TransactionNoteCol,
			TransactionDateCol,
		),
		query,
		query,
		limit,
	)
	if err != nil {
		return nil, err
	}
	result := make([]Transaction, 0, limit)
	for rows.Next() {
		tx := Transaction{}
		err := rows.Scan(&tx.Entity, &tx.Amount, &tx.Date, &tx.Note)
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	return result, nil
}

// InsertTransaction inserts a transaction into the transactions table
func (db *DB) InsertTransaction(tx Transaction) (sql.Result, error) {
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
func (db *DB) RemoveTransaction(tx Transaction) (sql.Result, error) {
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
func (db *DB) UpdateTransaction(old, new Transaction) (sql.Result, error) {
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

// Total returns the total of all the transactions in the database
// ? will this become slow over time?
func (db *DB) Total() (int, error) {
	row := db.QueryRow(
		fmt.Sprintf(
			"SELECT sum(%s) FROM %s",
			TransactionAmountCol,
			TransactionTableName,
		),
	)
	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("models: error querying for total: %w", err)
	}
	return total, nil
}
