package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/Anthony-Fiddes/budgeter/cli/budgeter"
	"github.com/Anthony-Fiddes/budgeter/internal/conf"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	_ "github.com/mattn/go-sqlite3"
)

func getDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get database path: %v\n", err)
	}
	return filepath.Join(home, ".budgeter.db")
}

func getConfigPath() string {
	config, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("could not get config path: %v\n", err)
	}
	return filepath.Join(config, "budgeter_config.json")
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
	configPath := getConfigPath()
	table := &transaction.Table{DB: db}
	err := table.Init()
	if err != nil {
		log.Fatalf("could not initialize database transactions table: %v\n", err)
	}
	app := budgeter.CLI{
		Config:       &conf.JSONFile{Path: configPath},
		DBPath:       dbPath,
		Transactions: &transaction.Table{DB: db},
	}
	os.Exit(app.Run(os.Args))
}
