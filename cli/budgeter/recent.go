package budgeter

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/month"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
	"github.com/cheynewallace/tabby"
)

type recent struct {
	limit        int
	search       string
	flip         bool
	month        string
	Transactions Table
}

func newRecent(c *CLI) *recent {
	result := recent{}
	result.Transactions = c.Transactions
	return &result
}

func (r recent) Name() string {
	return "recent"
}

//go:embed recentUsage.txt
var recentUsage string

func (r recent) Usage() string {
	return recentUsage
}

func multiParse(layouts []string, date string) (time.Time, error) {
	for _, layout := range layouts {
		result, err := time.Parse(layout, date)
		if err == nil {
			return result, nil
		}
	}
	return time.Time{}, fmt.Errorf("input date %q could not be parsed with any provided layout (%q)", date, layouts)
}

func (r recent) getTransactions() ([]transaction.Transaction, error) {
	if r.month != "" {
		fmts := []string{"January/06", "January"}
		inputMonth, err := multiParse(fmts, r.month)
		if inputMonth.Year() == 0 {
			now := time.Now()
			inputMonth = inputMonth.AddDate(now.Year(), 0, 0)
		}
		if err != nil {
			return nil, err
		}
		start := month.Start(inputMonth)
		end := month.End(inputMonth)
		rows, err := r.Transactions.Range(start, end, r.limit)
		if err != nil {
			return nil, err
		}
		transactions, err := rows.ScanSet()
		if err != nil {
			return nil, err
		}
		var result []transaction.Transaction
		for _, tx := range transactions {
			query := strings.ToLower(r.search)
			entity := strings.ToLower(tx.Entity)
			note := strings.ToLower(tx.Note)
			if !strings.Contains(entity, query) && !strings.Contains(note, query) {
				continue
			}
			result = append(result, tx)
		}
		return result, nil
	}

	rows, err := r.Transactions.Search(r.search, r.limit)
	if err != nil {
		return nil, err
	}
	transactions, err := rows.ScanSet()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// recent lists the most recently added transactions.
// TODO: Add a "pinned" feature/subcommand?
// TODO: Add a total for searches
func (r recent) Run(cmdArgs []string) error {
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

	fs := getFlagset(r.Name())
	fs.StringVar(&r.search, "s", "", "")
	fs.BoolVar(&r.flip, "f", false, "")
	fs.IntVar(&r.limit, "l", defaultRecentLimit, "")
	fs.StringVar(&r.month, "m", "", "")
	if err := fs.Parse(cmdArgs); err != nil {
		return err
	}
	fs.Usage()
	args := fs.Args()
	if len(args) > 0 {
		return fmt.Errorf("%s takes no arguments", r.Name())
	}

	tab := tabby.New()
	tab.AddHeader(idHeader, dateHeader, entityHeader, amountHeader, noteHeader)
	transactions, err := r.getTransactions()
	if err != nil {
		return err
	}
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
		monthTotal, err := r.Transactions.RangeTotal(month.Start(now), now)
		if err != nil {
			return err
		}
		totalStr := fmt.Sprintf("Current Month: %s", monthTotal)
		for i := 0; i < len(totalStr); i++ {
			fmt.Print("=")
		}
		fmt.Println()
		fmt.Println(totalStr)
	}
	return nil
}
