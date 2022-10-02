package budgeter

import (
	"fmt"
	"strconv"

	_ "embed"
)

type remove struct {
	Transactions Table
}

func newRemove(c *CLI) *remove {
	return &remove{Transactions: c.Transactions}
}

func (r remove) Name() string {
	return "remove"
}

//go:embed removeUsage.txt
var removeUsage string

func (r remove) Usage() string {
	return removeUsage
}

func (r remove) Run(cmdArgs []string) error {
	fs := getFlagset(r.Name())
	if err := fs.Parse(cmdArgs); err != nil {
		return err
	}
	args := fs.Args()
	if len(args) != 1 {
		return fmt.Errorf("%s takes one argument", r.Name())
	}
	txID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf(
			"%s takes a numerical ID. try `budgeter %s` to see some IDs.",
			r.Name(),
			recent{}.Name(),
		)
	}
	err = r.Transactions.Remove(txID)
	if err != nil {
		return fmt.Errorf("could not remove transaction #%d: %v", txID, err)
	}
	return nil
}
