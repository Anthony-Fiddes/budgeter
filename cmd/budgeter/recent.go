package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
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

type recentFlags struct {
	limit  int
	search string
	flip   bool
}

// recent lists the most recently added transactions.
// TODO: Show SQLite IDs so that I can reference transactions?
// otherwise maybe a hash?
// TODO: Add a "pinned" feature/subcommand?
func recent(table *transaction.Table, cmdArgs []string) error {
	var err error
	flags := recentFlags{}
	fs := flag.NewFlagSet(recentName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&flags.search, "s", "", "")
	fs.BoolVar(&flags.flip, "f", false, "")
	if err := fs.Parse(cmdArgs); err != nil {
		return err
	}
	args := fs.Args()
	if len(args) == 1 {
		flags.limit, err = strconv.Atoi(args[0])
		if err != nil {
			return errors.New("count must be a number")
		}
	} else {
		flags.limit = defaultRecentLimit
	}

	rows, err := table.Search(flags.search, flags.limit)
	if err != nil {
		return err
	}
	transactions, err := rows.ScanSet(flags.limit)
	if err != nil {
		return err
	}

	tab := tabby.New()
	tab.AddHeader(dateHeader, entityHeader, amountHeader, noteHeader)
	for i := 0; i < len(transactions); i++ {
		index := i
		if !flags.flip {
			index = len(transactions) - 1 - index
		}
		tx := transactions[index]
		// Align all the amount cells
		amount := tx.AmountString()
		if tx.Amount >= 0 {
			amount = " " + amount
		}
		tab.AddLine(tx.DateString(), tx.Entity, amount, tx.Note)
	}
	tab.Print()

	if flags.search == "" {
		total, err := table.Total()
		if err != nil {
			return err
		}
		totalString := fmt.Sprintf(totalTemplate, transaction.Dollars(total))
		for i := 0; i < len(totalString); i++ {
			fmt.Print("=")
		}
		fmt.Println()
		fmt.Print(totalString)
	}
	fmt.Println()
	return nil
}
