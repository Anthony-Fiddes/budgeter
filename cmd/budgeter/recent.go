package main

import (
	"flag"

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
		table.AddLine(t.DateString(), t.Entity, t.AmountString(), t.Note)
	}
	table.Print()
	return nil
}
