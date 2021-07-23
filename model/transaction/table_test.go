package transaction_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	_ "github.com/mattn/go-sqlite3"
)

func getMemTable() (*transaction.Table, error) {
	const URI = ":memory:"
	db, err := sql.Open("sqlite3", URI)
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
	defer table.DB.Close()

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
			t.Log(err)
			t.Fatalf("could not insert %+v into table", tx)
		}

		err = table.Insert(tx)
		if err == nil {
			t.Fatal("table is expected to error when inserting a transaction that already exists in the table")
		}
	}

	// Search Test
	for _, tx := range testData {
		rows, err := table.Search(tx.Entity, 1)
		if err != nil {
			t.Logf("transaction: %+v", tx)
			t.Fatalf(`Search failed: %v`, err)
		}

		rows.Next()
		result, err := rows.Scan()
		if err != nil || result != tx {
			t.Log(err)
			t.Fatalf(
				`Search %s did not return %+v but %+v`,
				tx.Entity, tx, result,
			)
		}
		if rows.Next() {
			t.Log("Rows.Next should be false after one call")
		}
	}

	// Total Test
	{
		sum := 0
		for _, tx := range testData {
			sum += tx.Amount
		}
		result, err := table.Total()
		if err != nil {
			t.Error(err)
		}
		if result != sum {
			t.Errorf("Total did not return %d but %d", sum, result)
		}
	}
}
