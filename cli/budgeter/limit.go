package budgeter

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const (
	limitName = "limit"
)

//go:embed limitUsage.txt
var limitUsage string

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
	_, err = transaction.Cents(args[0])
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
	lim := fmt.Sprintf("%s/%s", args[0], args[1])
	err = c.Config.Put(limitName, lim)
	if err != nil {
		c.Log.Printf("could not store limit: %v", err)
		return 1
	}

	return 0
}
