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

func getMemDb(t *testing.T) *DB {
	db, err := sql.Open(sqlite3, sqlite3URI)
	if err != nil {
		t.Fatalf("error creating an in memory database for testing: %s", err)
	}
	return &DB{db}
}

func TestCreateTransactionTable(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()

	_, err := db.CreateTransactionTable()
	if err != nil {
		t.Fatalf("error creating the transaction table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	_, err = db.Query(query)
	if err != nil {
		t.Fatalf(
			"error attempting to select from %s (which should now exist): %s",
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
	_, err := db.CreateTransactionTable()
	if err != nil {
		t.Fatalf("error creating the transaction table: %s", err)
	}

	_, err = db.InsertTransaction(testTX)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	_, err = db.InsertTransaction(testTX)
	if err == nil {
		t.Fatalf("error expected when inserting duplicate transactions: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	result, err := db.Query(query)
	if err != nil {
		t.Fatalf("error querying the test transaction from the table: %s", err)
	}
	next := result.Next()
	if !next {
		t.Fatalf("no row was returned when one was expected in the table.")
	}
	resultTX := Transaction{}
	result.Scan(&resultTX.Entity, &resultTX.Amount, &resultTX.Date, &resultTX.Note)
	if resultTX.Entity != testTX.Entity || resultTX.Amount != testTX.Amount {
		t.Fatalf("expected to receive an entry like %v, instead received %v", testTX, resultTX)
	}
}

func TestRemoveTransaction(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()
	testTX := Transaction{
		Entity: "Barack",
		Amount: 10,
	}
	_, err := db.CreateTransactionTable()
	if err != nil {
		t.Fatalf("error creating the transaction table: %s", err)
	}
	_, err = db.InsertTransaction(testTX)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	_, err = db.RemoveTransaction(testTX)
	if err != nil {
		t.Fatalf("error removing transaction from table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	result, err := db.Query(query)
	if result.Next() || err != nil {
		t.Fatal(
			"the table was supposed to be empty after removing its only transaction",
		)
	}
}

func TestGetTransactions(t *testing.T) {
	db := getMemDb(t)
	defer db.Close()
	testTX1 := Transaction{
		Entity: "Barack",
		Amount: 10,
		Date:   2,
	}
	testTX2 := testTX1
	testTX2.Date = 3
	_, err := db.CreateTransactionTable()
	if err != nil {
		t.Fatalf("error creating the transaction table: %s", err)
	}
	_, err = db.InsertTransaction(testTX1)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}
	_, err = db.InsertTransaction(testTX2)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	tt, err := db.GetTransactions(10)
	if err != nil {
		t.Fatalf("error getting transactions: %s", err)
	}

	if len(tt) != 2 {
		t.Fatal("not enough transactions were returned.")
	}
	if tt[0].Date != testTX2.Date {
		t.Error(
			"the returned transactions were not returned in order from most recent to least",
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
	_, err := db.CreateTransactionTable()
	if err != nil {
		t.Fatalf("error creating the transaction table: %s", err)
	}
	_, err = db.InsertTransaction(testTX)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	newTX := Transaction{Entity: "Obama", Amount: 20}
	_, err = db.UpdateTransaction(testTX, newTX)
	if err != nil {
		t.Fatalf("error updating transaction in table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", TransactionTableName)
	result, err := db.Query(query)
	if err != nil {
		t.Fatalf("error querying database: %s", err)
	}
	result.Next()
	resultTX := Transaction{}
	result.Scan(&resultTX.Entity, &resultTX.Amount, &resultTX.Date, &resultTX.Note)
	if resultTX.Entity != newTX.Entity || resultTX.Amount != newTX.Amount {
		t.Fatalf("expected to receive an entry like %v, instead received %v", newTX, resultTX)
	}
}
