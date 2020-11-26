package models

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlite3    = "sqlite3"
	sqlite3URI = ":memory:"
)

func getMemDb(t *testing.T) *sql.DB {
	db, err := sql.Open(sqlite3, sqlite3URI)
	if err != nil {
		t.Fatalf("Error creating an in memory database for testing: %s", err)
	}
	return db
}

func TestCreateTransactionTable(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()
	_, err := CreateTransactionTable(db)
	if err != nil {
		t.Fatalf("Error creating the transaction table: %s", err)
	}
	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	_, err = db.Query(query)
	if err != nil {
		t.Fatalf(
			"Error attempting to select from %s (which should now exist): %s",
			TransactionTableName,
			err,
		)
	}
}
