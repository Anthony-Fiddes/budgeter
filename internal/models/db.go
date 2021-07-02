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
	result, err := db.Exec(
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
	if err != nil {
		return nil, fmt.Errorf(
			"cannot create %v table: %w", TransactionTableName, err,
		)
	}
	return result, err
}

// TransactionRows wraps *sql.Rows to easily scan Transactions from a DB
type TransactionRows struct {
	*sql.Rows
}

func (tr *TransactionRows) Scan() (Transaction, error) {
	tx := Transaction{}
	err := tr.Rows.Scan(&tx.Entity, &tx.Amount, &tx.Date, &tx.Note)
	if err != nil {
		return Transaction{}, err
	}
	return tx, err
}

func (tr *TransactionRows) scanSet() ([]Transaction, error) {
	var result []Transaction
	for tr.Next() {
		tx, err := tr.Scan()
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	if err := tr.Err(); err != nil {
		return nil, fmt.Errorf("failed to read all transactions: %w", err)
	}
	return result, nil
}

// GetTransactions wraps GetTransactionRows to return a slice of Transactions. It
// cannot be used with a negative "limit".
func (db *DB) GetTransactions(limit int) ([]Transaction, error) {
	if limit < 1 {
		return nil, fmt.Errorf("cannot use negative limit %d", limit)
	}
	txRows, err := db.GetTransactionRows(limit)
	if err != nil {
		return nil, err
	}
	result, err := txRows.scanSet()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func queryError(e error) error {
	return fmt.Errorf("could not query database: %w", e)
}

// GetTransactionRows gets transactions from the transactions table. The
// most recent transactions will be returned first. GetTransactions will
// return, at most, "limit" results. A negative number can be passed to indicate
// that all possible results should be returned.
func (db *DB) GetTransactionRows(limit int) (*TransactionRows, error) {
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT * FROM %s ORDER BY %s DESC LIMIT ?",
			TransactionTableName,
			TransactionDateCol,
		),
		limit,
	)
	if err != nil {
		return nil, queryError(err)
	}
	return &TransactionRows{Rows: rows}, nil
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
		return nil, queryError(err)
	}
	txRows := TransactionRows{Rows: rows}
	result, err := txRows.scanSet()
	if err != nil {
		return nil, err
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
		return 0, fmt.Errorf("could not query database for total: %w", err)
	}
	return total, nil
}
