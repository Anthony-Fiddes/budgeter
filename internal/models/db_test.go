package models_test

import (
	"fmt"
	"testing"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
	"github.com/Anthony-Fiddes/budgeter/internal/modelstest"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateTransactionTable(t *testing.T) {
	db, err := modelstest.GetMemDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.CreateTransactionTable()
	if err != nil {
		t.Fatalf("error creating the transaction table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", models.TransactionTableName)
	_, err = db.Query(query)
	if err != nil {
		t.Fatalf(
			"error attempting to select from %s (which should now exist): %s",
			models.TransactionTableName,
			err,
		)
	}
}

func TestInsertTransaction(t *testing.T) {
	db, err := modelstest.GetMemDBWithTable()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testTX := models.Transaction{
		Entity: "bleep",
		Amount: 10,
	}

	_, err = db.InsertTransaction(testTX)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	_, err = db.InsertTransaction(testTX)
	if err == nil {
		t.Fatalf("error expected when inserting duplicate transactions: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", models.TransactionTableName)
	result, err := db.Query(query)
	if err != nil {
		t.Fatalf("error querying the test transaction from the table: %s", err)
	}
	next := result.Next()
	if !next {
		t.Fatalf("no row was returned when one was expected in the table.")
	}
	resultTX := models.Transaction{}
	result.Scan(&resultTX.Entity, &resultTX.Amount, &resultTX.Date, &resultTX.Note)
	if resultTX.Entity != testTX.Entity || resultTX.Amount != testTX.Amount {
		t.Fatalf("expected to receive an entry like %v, instead received %v", testTX, resultTX)
	}
}

func TestRemoveTransaction(t *testing.T) {
	db, err := modelstest.GetMemDBWithTable()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testTX := models.Transaction{
		Entity: "bleep",
		Amount: 10,
	}
	_, err = db.InsertTransaction(testTX)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	_, err = db.RemoveTransaction(testTX)
	if err != nil {
		t.Fatalf("error removing transaction from table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", models.TransactionTableName)
	result, err := db.Query(query)
	if result.Next() || err != nil {
		t.Fatal(
			"the table was supposed to be empty after removing its only transaction",
		)
	}
}

func TestGetTransactions(t *testing.T) {
	db, err := modelstest.GetMemDBWithTable()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testTX1 := models.Transaction{
		Entity: "bleep",
		Amount: 10,
		Date:   2,
	}
	testTX2 := testTX1
	testTX2.Date = 3
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

func TestSearch(t *testing.T) {
	db, err := modelstest.GetMemDBWithTable()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testTX1 := models.Transaction{
		Entity: "bleep",
		Amount: 10,
		Date:   2,
	}
	testTX2 := testTX1
	testTX2.Entity = "bloop"
	testTX2.Date = 3
	_, err = db.InsertTransaction(testTX1)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}
	_, err = db.InsertTransaction(testTX2)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	tt, err := db.Search("bl", 10)
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

	tt, err = db.Search("bloo", 10)
	if err != nil {
		t.Fatalf("error getting transactions: %s", err)
	}
	if len(tt) != 1 {
		t.Log(tt)
		t.Fatal("not enough transactions were returned.")
	}
}

func TestUpdateTransaction(t *testing.T) {
	db, err := modelstest.GetMemDBWithTable()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testTX := models.Transaction{
		Entity: "bleep",
		Amount: 10,
	}
	_, err = db.InsertTransaction(testTX)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}

	newTX := models.Transaction{Entity: "Obama", Amount: 20}
	_, err = db.UpdateTransaction(testTX, newTX)
	if err != nil {
		t.Fatalf("error updating transaction in table: %s", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s", models.TransactionTableName)
	result, err := db.Query(query)
	if err != nil {
		t.Fatalf("error querying database: %s", err)
	}
	result.Next()
	resultTX := models.Transaction{}
	result.Scan(&resultTX.Entity, &resultTX.Amount, &resultTX.Date, &resultTX.Note)
	if resultTX.Entity != newTX.Entity || resultTX.Amount != newTX.Amount {
		t.Fatalf("expected to receive an entry like %v, instead received %v", newTX, resultTX)
	}
}

func TestGetTotal(t *testing.T) {
	db, err := modelstest.GetMemDBWithTable()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	testTX1 := models.Transaction{
		Entity: "bleep",
		Amount: 10,
	}
	_, err = db.InsertTransaction(testTX1)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}
	total, err := db.Total()
	if err != nil {
		t.Fatalf("error getting table total: %s", err)
	}
	expected := testTX1.Amount
	if expected != total {
		t.Fatalf("expected a total of %d but got %d", expected, total)
	}
	testTX2 := models.Transaction{
		Entity: "bloop",
		Amount: 5,
	}
	_, err = db.InsertTransaction(testTX2)
	if err != nil {
		t.Fatalf("error inserting transaction into table: %s", err)
	}
	total, err = db.Total()
	if err != nil {
		t.Fatalf("error getting table total: %s", err)
	}
	expected = testTX1.Amount + testTX2.Amount
	if expected != total {
		t.Fatalf("expected a total of %d but got %d", expected, total)
	}
}
