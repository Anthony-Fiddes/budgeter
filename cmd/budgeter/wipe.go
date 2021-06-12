package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
)

const wipeName = "wipe"

func wipe(db *models.DB, cmdArgs []string) error {
	if len(cmdArgs) > 0 {
		return fmt.Errorf("%s does not take any arguments", wipeName)
	}
	// ? should this loop?
	fmt.Print("This will delete your budgeting information. Are you sure you want to continue? (y/[n]) ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	err := s.Err()
	if err != nil {
		return err
	}
	response := s.Text()
	response = strings.TrimSpace(response)
	response = strings.ToLower(response)
	if response != "y" {
		fmt.Println("No data deleted.")
		return nil
	}

	fmt.Println("Proceeding with deletion...")
	db.Close()
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}
	err = os.Remove(dbPath)
	if err != nil {
		return err
	}
	fmt.Println("Done.")
	return nil
}
