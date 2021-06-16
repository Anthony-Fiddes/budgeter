package main

import (
	"flag"
	"fmt"

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

func recent(db *models.DB, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(recentName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}
	// TODO: Add option to specify how many entries to print out
	transactions, err := db.GetTransactions(defaultRecentLimit)
	if err != nil {
		return err
	}
	table := tabby.New()
	table.AddHeader(dateHeader, entityHeader, amountHeader, noteHeader)
	for _, t := range transactions {
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
