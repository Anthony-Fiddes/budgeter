package budgeter

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"strconv"
	"time"

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
func recent(c *CLI) int {
	var err error
	flags := recentFlags{}
	fs := flag.NewFlagSet(recentName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&flags.search, "s", "", "")
	fs.BoolVar(&flags.flip, "f", false, "")
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		c.Log.Println()
		c.Log.Println(recentUsage)
		return 1
	}
	args := fs.Args()
	if len(args) == 1 {
		flags.limit, err = strconv.Atoi(args[0])
		if err != nil {
			c.Log.Printf("count \"%s\" must be a number", args[0])
			c.Log.Println()
			c.Log.Println(recentUsage)
			return 1
		}
	} else {
		flags.limit = defaultRecentLimit
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
		total, err := c.Transactions.Total()
		if err != nil {
			c.Log.Println(err)
			return 1
		}
		totalString := fmt.Sprintf(totalTemplate, transaction.Dollars(total))
		for i := 0; i < len(totalString); i++ {
			fmt.Print("=")
		}
		fmt.Println()
		fmt.Println(totalString)

		// TODO: make this configurable with limit subcommand
		// TODO: maybe add a test for this since it was buggy before?
		now := time.Now().UTC()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		rows, err := c.Transactions.Range(monthStart, now, -1)
		if err != nil {
			c.Log.Println(err)
			return 1
		}
		monthTxs, err := rows.ScanSet()
		if err != nil {
			c.Log.Println(err)
			return 1
		}
		oneMonthSpending := 0
		for _, tx := range monthTxs {
			oneMonthSpending += tx.Amount
		}
		fmt.Printf("Current Month: %s", transaction.Dollars(oneMonthSpending))
	}
	fmt.Println()
	return 0
}
