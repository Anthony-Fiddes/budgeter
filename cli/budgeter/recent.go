package budgeter

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/month"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	"github.com/cheynewallace/tabby"
)

const recentName = "recent"

//go:embed recentUsage.txt
var recentUsage string

// recent lists the most recently added transactions.
// TODO: Add a "pinned" feature/subcommand?
// TODO: Add a total for searches
func recent(c *CLI) int {
	type recentFlags struct {
		limit  int
		search string
		flip   bool
	}

	const (
		// defaultRecentLimit specifies the default number of items to receive when
		// the recent command is called
		defaultRecentLimit = 20
		idHeader           = "ID"
		dateHeader         = "Date"
		entityHeader       = "Entity"
		amountHeader       = "Amount"
		noteHeader         = "Note"
	)

	var err error
	flags := recentFlags{}
	fs := getFlagset(recentName)
	fs.StringVar(&flags.search, "s", "", "")
	fs.BoolVar(&flags.flip, "f", false, "")
	fs.IntVar(&flags.limit, "l", defaultRecentLimit, "")
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		c.Log.Println()
		c.Log.Println(recentUsage)
		return 1
	}
	fs.Usage()
	args := fs.Args()
	if len(args) > 0 {
		c.Log.Printf("%s takes no arguments", recentName)
		return 1
	}

	rows, err := c.Transactions.Search(flags.search, flags.limit)
	if err != nil {
		c.Log.Println(err)
		return 1
	}
	transactions, err := rows.ScanSet()
	if err != nil {
		c.Log.Println(err)
		return 1
	}

	tab := tabby.New()
	tab.AddHeader(idHeader, dateHeader, entityHeader, amountHeader, noteHeader)
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
		tab.AddLine(tx.ID, tx.DateString(), tx.Entity, amount, tx.Note)
	}
	tab.Print()

	if flags.search == "" {
		// TODO: make this configurable with limit subcommand
		// TODO: maybe add a test for this since it was buggy before?
		now := time.Now().UTC()
		monthTotal, err := c.Transactions.RangeTotal(month.Start(now), now)
		if err != nil {
			c.Log.Println(err)
			return 1
		}
		totalStr := fmt.Sprintf("Current Month: %s", transaction.Dollars(monthTotal))
		for i := 0; i < len(totalStr); i++ {
			fmt.Print("=")
		}
		fmt.Println()
		fmt.Println(totalStr)
	}
	return 0
}
