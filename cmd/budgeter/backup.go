package main

import (
	"fmt"
	"io"
	"os"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const backupName = "backup"

//go:embed backupUsage.txt
var backupUsage string

func backup(table *transaction.Table, cmdArgs []string) error {
	if len(cmdArgs) != 1 {
		return fmt.Errorf("%s only takes one argument", backupName)
	}
	table.DB.Close()
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}
	dbFile, err := os.Open(dbPath)
	if err != nil {
		return err
	}
	defer dbFile.Close()
	targetPath := cmdArgs[0]
	target, err := os.OpenFile(targetPath, os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	_, err = io.Copy(target, dbFile)
	if err != nil {
		return err
	}
	return nil
}
