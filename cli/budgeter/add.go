package budgeter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

var addUsage = "add doesn't quite have options just yet!"

type add struct {
	*CLI
	lastDate string
	lastUnix int64
}

func (a *add) Name() string {
	return "add"
}

func (a *add) Run(c *CLI) int {
	a.CLI = c
	fs := getFlagset(a.Name())
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		c.err.Println()
		c.err.Print(addUsage)
		return 1
	}
	args := fs.Args()
	if len(args) == 0 {
		return a.interactiveAdd()
	} else if len(args) > fieldsPerRecord {
		c.err.Printf("%s takes at most %d arguments", a.Name(), fieldsPerRecord)
		c.err.Println()
		c.err.Print(addUsage)
		return 1
	}
	// TODO: implement an option that parses from flags or from args
	return 0
}

// TODO: Find a way to handle duplicates gracefully
func (a *add) interactiveAdd() int {
	// TODO: allow short dates like "21" or "6/21" that
	// default to this month or year
	for {
		tx, err := a.getTransaction()
		if err != nil {
			a.err.Println(err)
			return 1
		}
		if err := a.Transactions.Insert(tx); err != nil {
			a.err.Println(err)
			return 1
		}

		// TODO: Add context when adding transactions between sessions.
		// e.g. making the last used date the new default?, enabling an undo command
		fmt.Fprint(a.Out, "\nWould you like to add another transaction? (y/[n]) ")
		confirmed, err := inpt.Confirm()
		fmt.Println()
		if err != nil {
			a.err.Println(err)
			return 1
		}
		if !confirmed {
			break
		}
	}

	return 0
}

func (a *add) getField(field string) (string, error) {
	fmt.Fprintf(a.Out, "%s: ", field)
	response, err := a.in.Line()
	if err != nil {
		return "", err
	}
	return response, err
}

func (a *add) getDate() (int64, error) {
	if a.lastDate == "" {
		now := time.Now()
		a.lastDate = now.Format(transaction.DateLayout)
		a.lastUnix = now.Unix()
	}
	fmt.Fprintf(a.Out, "%s [%s]: ", transaction.DateCol, a.lastDate)
	response, err := a.in.Line()
	if err != nil {
		return 0, err
	}
	if response == "" {
		// the default date is the last date the user entered. If they haven't
		// entered anything, it's today's date.
		response = a.lastDate
	}
	if strings.Count(response, "/") == 1 {
		// the user may enter a short date in the format of MM/DD
		// the year will be taken from their last response.
		response += "/"
		t := time.Unix(a.lastUnix, 0)
		lastYear := strconv.Itoa(t.Year())
		response += lastYear
	}
	unix, err := transaction.Unix(response)
	if err != nil {
		return 0, err
	}
	a.lastDate = response
	a.lastUnix = unix
	return unix, err
}

func (a *add) getTransaction() (transaction.Transaction, error) {
	var err error
	tx := transaction.Transaction{}
	tx.Date, err = a.getDate()
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Entity, err = a.getField(transaction.EntityCol)
	if err != nil {
		return transaction.Transaction{}, err
	}
	amount, err := a.getField(transaction.AmountCol)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Amount, err = transaction.GetCents(amount)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Note, err = a.getField(transaction.NoteCol)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return tx, nil
}
