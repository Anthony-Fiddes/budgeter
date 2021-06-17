package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
)

const (
	wipeName          = "wipe"
	wipeCancelMessage = "No data deleted."
)

func wipe(db *models.DB, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(wipeName, flag.ContinueOnError)
	confirmed := fs.Bool("y", false, "Confirms that the user would like to wipe their budgeting information.")
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		fs.Usage()
		return fmt.Errorf("%s does not take any arguments", wipeName)
	}
	if *confirmed {
		return wipeDB(db)
	}
	return interactiveWipe(db)
}

func wipeDB(db *models.DB) error {
	db.Close()
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}
	err = os.Remove(dbPath)
	if err != nil {
		return err
	}
	fmt.Println("Done. All budgeting information deleted.")
	return nil
}

func interactiveWipe(db *models.DB) error {
	// ? should this loop?
	fmt.Print("This will delete your budgeting information. Are you sure you want to continue? (y/[n]) ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return err
	}
	response = strings.TrimSpace(response)
	response = strings.ToLower(response)
	if response != "y" {
		fmt.Println(wipeCancelMessage)
		return nil
	}
	fmt.Println("Proceeding with deletion...")
	return wipeDB(db)
}
