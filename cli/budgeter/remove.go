package budgeter

import (
	"strconv"

	_ "embed"
)

const (
	removeName = "remove"
)

//go:embed removeUsage.txt
var removeUsage string

func remove(c *CLI) int {
	fs := getFlagset(removeName)
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		c.Err.Println()
		c.Err.Print(removeUsage)
		return 1
	}
	args := fs.Args()
	if len(args) != 1 {
		c.Err.Printf("%s takes one argument", removeName)
		c.Err.Println()
		c.Err.Print(removeUsage)
		return 1
	}
	txID, err := strconv.Atoi(args[0])
	if err != nil {
		c.Err.Printf(
			"%s takes a numerical ID. try `budgeter %s` to see some IDs.",
			removeName,
			recentName,
		)
		c.Err.Println()
		c.Err.Print(removeUsage)
		return 1
	}
	err = c.Transactions.Remove(txID)
	if err != nil {
		c.Err.Printf("could not remove transaction #%d: %v", txID, err)
		return 1
	}
	return 0
}
