package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed usage.txt
var usage string

type command func(*models.DB, []string) error

const dbName = "budgeter.db"

func getDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".", fmt.Errorf("error getting db path: %w", err)
	}
	return filepath.Join(home, dbName), nil
}

func initDB() (*models.DB, error) {

	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}
	d, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db := &models.DB{DB: d}
	_, err = db.CreateTransactionTable()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func printUsage() {
	fmt.Fprint(os.Stderr, usage)
	os.Exit(1)
}

func main() {
	commands := map[string]command{
		ingestName: ingest,
		wipeName:   wipe,
		recentName: recent,
	}
	if len(os.Args) < 2 {
		printUsage()
	}
	alias := os.Args[1]
	args := os.Args[2:]

	db, err := initDB()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()

	cmd, ok := commands[alias]
	if !ok {
		printUsage()
	}
	err = cmd(db, args)
	if err != nil {
		log.Fatalf("%v: %v", alias, err)
	}
}
