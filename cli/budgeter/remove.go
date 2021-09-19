package budgeter

import (
	"flag"
	"strconv"

	_ "embed"
)

const (
	removeName = "remove"
)

//go:embed removeUsage.txt
var removeUsage string

func remove(c *CLI) int {
	fs := flag.NewFlagSet(removeName, flag.ContinueOnError)
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		return 1
	}
	args := fs.Args()
	if len(args) != 1 {
		c.Log.Printf("%s takes one argument", removeName)
		c.Log.Println()
		c.Log.Print(removeUsage)
		return 1
	}
	txID, err := strconv.Atoi(args[0])
	if err != nil {
		c.Log.Printf(
			"%s takes a numerical ID. try `budgeter %s` to see some IDs.",
			removeName,
			recentName,
		)
		c.Log.Println()
		c.Log.Print(removeUsage)
		return 1
	}
	err = c.Transactions.Remove(txID)
	if err != nil {
		c.Log.Printf("could not remove transaction #%d: %v", txID, err)
		return 1
	}
	return 0
}
