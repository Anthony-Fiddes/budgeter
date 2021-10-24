package budgeter

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/month"
	"github.com/cheynewallace/tabby"
)

type recent struct {
	limit  int
	search string
	flip   bool
}

//go:embed recentUsage.txt
var recentUsage string

// recentName is a const because it is used as a suggestion in remove.
const recentName = "recent"

func (r *recent) Name() string {
	return recentName
}

// recent lists the most recently added transactions.
// TODO: Add a "pinned" feature/subcommand?
// TODO: Add a total for searches
func (r *recent) Run(c *CLI) int {
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
	fs := getFlagset(r.Name())
	fs.StringVar(&r.search, "s", "", "")
	fs.BoolVar(&r.flip, "f", false, "")
	fs.IntVar(&r.limit, "l", defaultRecentLimit, "")
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		c.err.Println()
		c.err.Println(recentUsage)
		return 1
	}
	fs.Usage()
	args := fs.Args()
	if len(args) > 0 {
		c.err.Printf("%s takes no arguments", r.Name())
		return 1
	}

	rows, err := c.Transactions.Search(r.search, r.limit)
	if err != nil {
		c.err.Println(err)
		return 1
	}
	transactions, err := rows.ScanSet()
	if err != nil {
		c.err.Println(err)
		return 1
	}

	tab := tabby.New()
	tab.AddHeader(idHeader, dateHeader, entityHeader, amountHeader, noteHeader)
	for i := 0; i < len(transactions); i++ {
		index := i
		if !r.flip {
			index = len(transactions) - 1 - index
		}
		tx := transactions[index]
		// Align all the amount cells
		amount := tx.Amount.String()
		if tx.Amount >= 0 {
			amount = " " + amount
		}
		tab.AddLine(tx.ID, tx.DateString(), tx.Entity, amount, tx.Note)
	}
	tab.Print()

	if r.search == "" {
		// TODO: make this configurable with limit subcommand
		// TODO: maybe add a test for this since it was buggy before?
		now := time.Now().UTC()
		monthTotal, err := c.Transactions.RangeTotal(month.Start(now), now)
		if err != nil {
			c.err.Println(err)
			return 1
		}
		totalStr := fmt.Sprintf("Current Month: %s", monthTotal)
		for i := 0; i < len(totalStr); i++ {
			fmt.Print("=")
		}
		fmt.Println()
		fmt.Println(totalStr)
	}
	return 0
}
