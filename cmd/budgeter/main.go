package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed usage.txt
var usage string

type command struct {
	Exec  func(*transaction.Table, []string) error
	Usage string
}

const dbName = "budgeter.db"

func getDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".", fmt.Errorf("error getting database path: %w", err)
	}
	return filepath.Join(home, dbName), nil
}

func initDB() (*sql.DB, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return db, nil
}

func stderr(usage string) {
	fmt.Fprint(os.Stderr, usage)
}

func main() {
	// TODO: Add an export command
	// TODO: Add a remove/edit command
	// TODO: Add a query command
	// TODO: Consider adding a period command to search through a certain time period
	commands := map[string]command{
		addName:    {Exec: add},
		backupName: {Exec: backup, Usage: backupUsage},
		ingestName: {Exec: ingest},
		wipeName:   {Exec: wipe},
		recentName: {Exec: recent, Usage: recentUsage},
	}
	if len(os.Args) < 2 {
		stderr(usage)
		os.Exit(1)
	}
	alias := os.Args[1]
	args := os.Args[2:]

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	table := &transaction.Table{DB: db}
	table.Init()

	cmd, ok := commands[alias]
	if !ok {
		stderr(usage)
		os.Exit(1)
	}
	err = cmd.Exec(table, args)
	if err != nil {
		log.Printf("%v: %v\n", alias, err)
		if cmd.Usage != "" {
			stderr("\n")
			stderr(cmd.Usage)
		}
		os.Exit(1)
	}
}
