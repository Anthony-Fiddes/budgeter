package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/Anthony-Fiddes/budgeter/cli/budgeter"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	_ "github.com/mattn/go-sqlite3"
)

const dbName = ".budgeter.db"

func getDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get database path: %v\n", err)
	}
	return filepath.Join(home, dbName)
}

func initDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	return db
}

func main() {
	dbPath := getDBPath()
	db := initDB(dbPath)
	table := &transaction.Table{DB: db}
	err := table.Init()
	if err != nil {
		log.Fatalf("could not initialize database transactions table: %v\n", err)
	}
	app := budgeter.CLI{
		DBPath:       dbPath,
		Transactions: &transaction.Table{DB: db},
	}
	os.Exit(app.Run(os.Args))
}
