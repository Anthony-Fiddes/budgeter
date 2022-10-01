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
	cache, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("could not get database path: %v\n", err)
	}
	return filepath.Join(cache, "budgeter.tsv")
}

func getConfigPath() string {
	config, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("could not get config path: %v\n", err)
	}
	config = filepath.Join(config, "budgeter")
	err = os.MkdirAll(config, 0755)
	if err != nil {
		log.Fatalf("could not create the config directory: %v\n", err)
	}
	return filepath.Join(config, "config.json")
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
	configPath := getConfigPath()
	app := budgeter.CLI{
		Config:       conf.NewJSONFile(configPath),
		DBPath:       dbPath,
		Transactions: transaction.NewTSVTable(getDBPath()),
	}
	os.Exit(app.Run(os.Args))
}
