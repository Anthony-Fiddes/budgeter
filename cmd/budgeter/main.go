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

type command struct {
	Exec  func(*models.DB, []string) error
	Usage string
}

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

func stderr(usage string) {
	fmt.Fprint(os.Stderr, usage)
}

func main() {
	// TODO: Add a backup command
	// TODO: Add a query command
	// TODO: Add a search command
	// TODO: Add a remove command
	// TODO: Add an edit command
	// TODO: Consider adding a period command to search through a certain time period
	commands := map[string]command{
		addName:    {Exec: add},
		ingestName: {Exec: ingest},
		wipeName:   {Exec: wipe},
		recentName: {Exec: recent, Usage: recentUsage},
		backupName: {Exec: backup, Usage: backupUsage},
	}
	if len(os.Args) < 2 {
		stderr(usage)
		os.Exit(1)
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
		stderr(usage)
		os.Exit(1)
	}
	err = cmd.Exec(db, args)
	if err != nil {
		log.Printf("%v: %v\n", alias, err)
		if cmd.Usage != "" {
			stderr("\n")
			stderr(cmd.Usage)
		}
	}
}
