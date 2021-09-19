package transaction

import (
	"database/sql"
	"fmt"
	"time"
)

// Table is the transactions table in a database
type Table struct{ DB *sql.DB }

// Init creates the transactions table if it doesn't exist.
func (t *Table) Init() error {
	_, err := t.DB.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s "+
				"(%s INTEGER NOT NULL PRIMARY KEY, "+
				"%s TEXT NOT NULL, %s INTEGER NOT NULL, %s INTEGER NOT NULL, %s TEXT NOT NULL, "+
				"UNIQUE(%s,%s,%s,%s))",
			TableName,
			IDCol,
			EntityCol,
			AmountCol,
			DateCol,
			NoteCol,
			EntityCol,
			AmountCol,
			DateCol,
			NoteCol,
		),
	)
	if err != nil {
		return fmt.Errorf(
			"transaction: cannot create table: %w", err,
		)
	}
	return nil
}

func queryError(e error) error {
	return fmt.Errorf("transaction: could not query table: %w", e)
}

// Search returns the most recent transactions that include the given "query".
// It returns, at most, "limit" transactions, and returns more recent
// transactions first. A negative "limit" will return as many
// transactions as are available.
func (t *Table) Search(query string, limit int) (*Rows, error) {
	query = "%" + query + "%"
	rows, err := t.DB.Query(
		fmt.Sprintf(
			"SELECT * FROM %s WHERE %s LIKE ? OR %s LIKE ? ORDER BY %s DESC LIMIT ?",
			TableName,
			EntityCol,
			NoteCol,
			DateCol,
		),
		query,
		query,
		limit,
	)
	if err != nil {
		return nil, queryError(err)
	}
	return &Rows{rows}, nil
}

// Range returns the transactions that occurred within the give range of time.
// It returns, at most, "limit" transactions, and returns them in chronological
// order. A negative "limit" will return as many transactions as are available.
func (t *Table) Range(start, end time.Time, limit int) (*Rows, error) {
	startUnix := start.UTC().Unix()
	stopUnix := end.UTC().Unix()
	rows, err := t.DB.Query(
		fmt.Sprintf(
			"SELECT * FROM %s WHERE %s >= ? AND %s <= ? ORDER BY %s ASC LIMIT ?",
			TableName,
			DateCol,
			DateCol,
			DateCol,
		),
		startUnix,
		stopUnix,
		limit,
	)
	if err != nil {
		return nil, queryError(err)
	}
	return &Rows{rows}, nil
}

// RangeTotal returns the cost of the transactions that occurred within the give
// range of time.
//
// It uses, at most, "limit" transactions. A negative "limit" will use as many
// transactions as are available.
func (t *Table) RangeTotal(start, end time.Time) (int, error) {
	startUnix := start.UTC().Unix()
	stopUnix := end.UTC().Unix()
	row := t.DB.QueryRow(
		fmt.Sprintf(
			"SELECT SUM(%s) FROM %s WHERE %s >= ? AND %s <= ? ORDER BY %s ASC",
			AmountCol,
			TableName,
			DateCol,
			DateCol,
			DateCol,
		),
		startUnix,
		stopUnix,
	)
	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("could not get total: %w", err)
	}
	return total, nil
}

// Insert inserts a transaction into the transactions table. The ID provided by
// "tx" is ignored, as the database determines the ID.
func (t *Table) Insert(tx Transaction) error {
	_, err := t.DB.Exec(
		fmt.Sprintf(
			"INSERT INTO %s(%s, %s, %s, %s) VALUES (?, ?, ?, ?)",
			TableName,
			EntityCol,
			AmountCol,
			DateCol,
			NoteCol,
		),
		tx.Entity,
		tx.Amount,
		tx.Date,
		tx.Note,
	)
	if err != nil {
		return fmt.Errorf("transaction: could not insert %+v: %w", tx, err)
	}
	return nil
}

// Total returns the total of all the transactions in the database
// ? will this become slow over time?
func (t *Table) Total() (int, error) {
	row := t.DB.QueryRow(
		fmt.Sprintf(
			"SELECT sum(%s) FROM %s",
			AmountCol,
			TableName,
		),
	)
	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("transaction: could not query database for total: %w", err)
	}
	return total, nil
}

// Remove deletes the given transaction from the table.
func (t *Table) Remove(transactionID int) error {
	_, err := t.DB.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE %s=?",
			TableName,
			IDCol,
		),
		transactionID,
	)
	if err != nil {
		return fmt.Errorf(
			"transaction: could not remove transaction #%d: %w",
			transactionID,
			err,
		)
	}
	return nil
}

// Rows wraps *sql.Rows to easily scan Transactions from a DB
type Rows struct{ *sql.Rows }

// Scan scans a transaction from the current result set.
func (r *Rows) Scan() (Transaction, error) {
	tx := Transaction{}
	err := r.Rows.Scan(&tx.ID, &tx.Entity, &tx.Amount, &tx.Date, &tx.Note)
	if err != nil {
		return Transaction{}, err
	}
	return tx, err
}

// ScanSet scans up to "limit" transactions from a result set
// into a slice. Do not use ScanSet if you expect that that your result set will
// be very large.
func (r *Rows) ScanSet() ([]Transaction, error) {
	var result []Transaction
	for r.Next() {
		tx, err := r.Scan()
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("transaction: failed to scan result set: %w", err)
	}
	return result, nil
}
