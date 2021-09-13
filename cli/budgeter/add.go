package budgeter

import (
	"flag"
	"fmt"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const (
	addName = "add"
)

// TODO: Find a way to handle duplicates gracefully
func interactiveAdd(c *CLI) int {
	getField := func(field string) (string, error) {
		fmt.Printf("%s: ", field)
		response, err := inpt.Line()
		if err != nil {
			return "", err
		}
		return response, err
	}

	lastDate := time.Now().Format(transaction.DateLayout)
	getDate := func() (int64, error) {
		fmt.Printf("%s [%s]: ", transaction.DateCol, lastDate)
		response, err := inpt.Line()
		if err != nil {
			return 0, err
		}
		if response == "" {
			response = lastDate
		}
		date, err := transaction.Date(response)
		if err != nil {
			return 0, err
		}
		lastDate = response
		return date, err
	}

	// TODO: allow short dates like "21" or "6/21" that
	// default to this month or year
	getTransaction := func() (transaction.Transaction, error) {
		var err error
		tx := transaction.Transaction{}
		tx.Date, err = getDate()
		if err != nil {
			return transaction.Transaction{}, err
		}
		tx.Entity, err = getField(transaction.EntityCol)
		if err != nil {
			return transaction.Transaction{}, err
		}
		amount, err := getField(transaction.AmountCol)
		if err != nil {
			return transaction.Transaction{}, err
		}
		tx.Amount, err = transaction.Cents(amount)
		if err != nil {
			return transaction.Transaction{}, err
		}
		tx.Note, err = getField(transaction.NoteCol)
		if err != nil {
			return transaction.Transaction{}, err
		}
		return tx, nil
	}

	for {
		tx, err := getTransaction()
		if err != nil {
			c.Log.Println(err)
			return 1
		}
		if err := c.Transactions.Insert(tx); err != nil {
			c.Log.Println(err)
			return 1
		}

		// TODO: Add context when adding transactions. e.g. making the last
		// used date the new default?, enabling an undo command
		fmt.Print("\nWould you like to add another transaction? (y/[n]) ")
		confirmed, err := inpt.Confirm()
		fmt.Println()
		if err != nil {
			c.Log.Println(err)
			return 1
		}
		if !confirmed {
			break
		}
	}

	return 0
}

func add(c *CLI) int {
	fs := flag.NewFlagSet(addName, flag.ContinueOnError)
	if err := fs.Parse(c.args); err != nil {
		c.logParsingErr(err)
		return 1
	}
	args := fs.Args()
	if len(args) == 0 {
		return interactiveAdd(c)
	} else if len(args) > fieldsPerRecord {
		c.Log.Printf("%s takes at most %d arguments", addName, fieldsPerRecord)
		return 1
	}
	// TODO: implement an option that parses from flags or from args
	return 0
}
