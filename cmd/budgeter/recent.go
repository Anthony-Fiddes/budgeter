package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
	"github.com/cheynewallace/tabby"
)

const (
	recentName = "recent"
	// defaultRecentLimit specifies the default number of items to receive when
	// the command is called
	defaultRecentLimit = 5
	dateHeader         = "Date"
	entityHeader       = "Entity"
	amountHeader       = "Amount"
	noteHeader         = "Note"
	totalTemplate      = "Total: %s"
)

//go:embed recentUsage.txt
var recentUsage string

// recent lists the most recently added transactions.
// TODO: Show SQLite IDs so that I can reference transactions?
// otherwise maybe a hash?
// TODO: Add a "pinned" feature/subcommand
func recent(db *models.DB, cmdArgs []string) error {
	recentLimit := defaultRecentLimit
	var err error
	fs := flag.NewFlagSet(recentName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}
	args := fs.Args()
	if len(args) == 1 {
		recentLimit, err = strconv.Atoi(args[0])
		if err != nil {
			return errors.New("count must be a number")
		}
	}
	transactions, err := db.GetTransactions(recentLimit)
	if err != nil {
		return err
	}
	table := tabby.New()
	table.AddHeader(dateHeader, entityHeader, amountHeader, noteHeader)
	// Reverse the order to display the most recent transactions at the bottom
	for i := len(transactions) - 1; i >= 0; i-- {
		t := transactions[i]
		// Align all the amount cells
		amount := t.AmountString()
		if t.Amount >= 0 {
			amount = " " + amount
		}
		table.AddLine(t.DateString(), t.Entity, amount, t.Note)
	}
	total, err := db.Total()
	if err != nil {
		return err
	}
	table.Print()
	totalString := fmt.Sprintf(totalTemplate, models.Dollars(total))
	for i := 0; i < len(totalString); i++ {
		fmt.Print("=")
	}
	fmt.Println()
	fmt.Print(totalString)
	fmt.Println()
	return nil
}
