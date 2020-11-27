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

func TestInsertTransaction(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()
	testTX := Transaction{
		Entity: "Barack",
		Amount: 10,
	}
	// ? Is there a way to remove this dependency? Is it fine as is?
	_, err := CreateTransactionTable(db)
	if err != nil {
		t.Fatalf("Error creating the transaction table: %s", err)
	}

	_, err = InsertTransaction(db, testTX)
	if err != nil {
		t.Fatalf("Error inserting transaction into table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	result, err := db.Query(query)
	if err != nil {
		t.Fatalf("Error querying the test transaction from the table: %s", err)
	}
	next := result.Next()
	if !next {
		t.Fatalf("No row was returned when one was expected in the table.")
	}
	resultTX := Transaction{}
	result.Scan(&resultTX.Entity, &resultTX.Amount, &resultTX.Date, &resultTX.Note)
	if resultTX.Entity != testTX.Entity || resultTX.Amount != testTX.Amount {
		t.Fatalf("Expected to receive an entry like %v, instead received %v", testTX, resultTX)
	}
}

func TestRemoveTransaction(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()
	testTX := Transaction{
		Entity: "Barack",
		Amount: 10,
	}
	_, err := CreateTransactionTable(db)
	if err != nil {
		t.Fatalf("Error creating the transaction table: %s", err)
	}
	_, err = InsertTransaction(db, testTX)
	if err != nil {
		t.Fatalf("Error inserting transaction into table: %s", err)
	}

	_, err = RemoveTransaction(db, testTX)
	if err != nil {
		t.Fatalf("Error removing transaction from table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	result, err := db.Query(query)
	if result.Next() {
		t.Fatal(
			"The table was supposed to be empty after removing its only transaction",
		)
	}
}

func TestUpdateTransaction(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()
	testTX := Transaction{
		Entity: "Barack",
		Amount: 10,
	}
	_, err := CreateTransactionTable(db)
	if err != nil {
		t.Fatalf("Error creating the transaction table: %s", err)
	}
	_, err = InsertTransaction(db, testTX)
	if err != nil {
		t.Fatalf("Error inserting transaction into table: %s", err)
	}

	newTX := Transaction{Entity: "Obama", Amount: 20}
	_, err = UpdateTransaction(db, testTX, newTX)
	if err != nil {
		t.Fatalf("Error updating transaction in table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	result, err := db.Query(query)
	result.Next()
	resultTX := Transaction{}
	result.Scan(&resultTX.Entity, &resultTX.Amount, &resultTX.Date, &resultTX.Note)
	if resultTX.Entity != newTX.Entity || resultTX.Amount != newTX.Amount {
		t.Fatalf("Expected to receive an entry like %v, instead received %v", newTX, resultTX)
	}
}
