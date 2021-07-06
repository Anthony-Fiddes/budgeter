package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const (
	wipeName          = "wipe"
	wipeCancelMessage = "No data deleted."
)

func wipe(table *transaction.Table, cmdArgs []string) error {
	fs := flag.NewFlagSet(wipeName, flag.ContinueOnError)
	table.DB.Close()
	confirmed := fs.Bool("y", false, "Confirms that the user would like to wipe their budgeting information.")
	err := fs.Parse(cmdArgs)
	if err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		fs.Usage()
		return fmt.Errorf("%s does not take any arguments", wipeName)
	}
	if *confirmed {
		return wipeDB()
	}
	return interactiveWipe()
}

func wipeDB() error {
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}
	if err := os.Remove(dbPath); err != nil {
		return err
	}
	fmt.Println("Done. All budgeting information deleted.")
	return nil
}

func interactiveWipe() error {
	// ? should this loop?
	fmt.Print("This will delete your budgeting information. Are you sure you want to continue? (y/[n]) ")
	confirmed, err := inpt.Confirm()
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}
	fmt.Println("Proceeding with deletion...")
	return wipeDB()
}
