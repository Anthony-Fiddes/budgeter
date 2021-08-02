package budgeter

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const (
	budgetKey = "budget"
	budgetSep = "/"
	limitName = "limit"
)

//go:embed limitUsage.txt
var limitUsage string

// setBudget stores the budget amount and period in a human readable format in the app's
// config store.
func (c *CLI) setBudget(cents int, p period) error {
	budget := fmt.Sprintf("%s%s%s", transaction.Dollars(cents), budgetSep, p.String())
	err := c.Config.Put(budgetKey, budget)
	if err != nil {
		return fmt.Errorf("could not store budget: %v", err)
	}
	return nil
}

// getBudget returns the user's specified budgeting limit in cents per period.
//
// e.g. 10000 cents / week
func (c *CLI) getBudget() (int, period, error) {
	budgetStr, err := c.Config.Get(budgetKey)
	if err != nil {
		return 0, unknown, fmt.Errorf("could not get budget: %w", err)
	}
	budget := strings.SplitN(budgetStr, budgetSep, 1)
	if len(budget) != 2 {
		return 0, unknown, fmt.Errorf("budget in store (%s) is formatted improperly: %w", budgetStr, err)
	}
	cents, err := transaction.Cents(budget[0])
	if err != nil {
		return 0, unknown, err
	}
	period := getPeriod(budget[1])
	if period.Unknown() {
		return 0, unknown, fmt.Errorf("unknown period \"%s\"", budget[1])
	}
	return cents, period, nil
}

func limit(c *CLI) int {
	fs := flag.NewFlagSet(exportName, flag.ContinueOnError)
	err := fs.Parse(c.args)
	if err != nil {
		c.logParsingErr(err)
		c.Log.Println()
		c.Log.Println(limitUsage)
		return 1
	}
	args := fs.Args()
	if len(args) != 2 {
		c.Log.Printf("%s takes two arguments", exportName)
		c.Log.Println()
		c.Log.Println(limitUsage)
		return 1
	}

	// Check that the user input is valid
	args[0] = inpt.Normalize(args[0])
	amount, err := transaction.Cents(args[0])
	if err != nil {
		c.Log.Println(err)
		return 1
	}
	args[1] = inpt.Normalize(args[1])
	per := getPeriod(args[1])
	if per == unknown {
		c.Log.Printf("invalid period \"%s\"", args[1])
		c.Log.Println()
		c.Log.Println(limitUsage)
		return 1
	}

	// Store the limit amount and period in a human readable format in the app's
	// config store.
	err = c.setBudget(amount, per)
	if err != nil {
		c.Log.Println(err)
		return 1
	}

	return 0
}
