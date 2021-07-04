package transaction_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	_ "github.com/mattn/go-sqlite3"
)

func getMemTable() (*transaction.Table, error) {
	const sqlite3URI = ":memory:"
	db, err := sql.Open("sqlite3", sqlite3URI)
	if err != nil {
		return nil, fmt.Errorf("error creating an in-memory database for testing: %w", err)
	}
	table := &transaction.Table{DB: db}
	err = table.Init()
	if err != nil {
		return nil, fmt.Errorf("error creating the transaction table: %w", err)
	}
	return table, nil
}

// TestTable tests Table and its methods all at once since they're all very coupled.
func TestTable(t *testing.T) {
	table, err := getMemTable()
	if err != nil {
		t.Fatal(err)
	}

	testData := []transaction.Transaction{
		// Test data graciously provided by Sarah Werum
		{
			Entity: "Apossumtheosis",
			Amount: 400000,
			Date:   -1,
			Note:   "it has begun.",
		},
		{
			Entity: "Squirrel Sanctuary",
			Amount: 123400,
			Date:   0,
			Note:   "Squirrels are very aggresive",
		},
		{
			Entity: "Frog rebellion",
			Amount: 30800,
			Date:   5,
			Note:   "they have risen up!",
		},
		// My test data
		{
			Entity: "Lyft",
			Amount: 1368,
			Date:   6,
			Note:   "Ride to the doctor",
		},
		{
			Entity: "Kroger",
			Amount: 1212,
			Date:   6,
			Note:   "Groceries",
		},
	}

	// Insert test
	for _, tx := range testData {
		err := table.Insert(tx)
		if err != nil {
			t.Fatalf("could not insert %+v into table", tx)
		}

		err = table.Insert(tx)
		if err == nil {
			const reason = "table is expected to error when inserting a transaction that already exists in the table"
			t.Fatal(reason)
		}
	}
}
