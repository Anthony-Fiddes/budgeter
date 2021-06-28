package main

import (
	"flag"
	"fmt"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/internal/models"
)

const (
	addName = "add"
)

// TODO: Find a way to handle duplicates gracefully
func interactiveAdd(db *models.DB) error {
	getField := func(field string) (string, error) {
		fmt.Printf("%s: ", field)
		response, err := inpt.Line()
		if err != nil {
			return "", err
		}
		return response, err
	}

	// TODO: make the date default to today
	// TODO: perhaps implement a date picker? / start the next date where you
	// left off with the last one / allow short dates like "21" or "6/21" that
	// default to this month or year
	getTransaction := func() (models.Transaction, error) {
		tx := models.Transaction{}
		date, err := getField(models.TransactionDateCol)
		if err != nil {
			return models.Transaction{}, err
		}
		tx.Date, err = models.Date(date)
		if err != nil {
			return models.Transaction{}, err
		}
		tx.Entity, err = getField(models.TransactionEntityCol)
		if err != nil {
			return models.Transaction{}, err
		}
		amount, err := getField(models.TransactionAmountCol)
		if err != nil {
			return models.Transaction{}, err
		}
		tx.Amount, err = models.Cents(amount)
		if err != nil {
			return models.Transaction{}, err
		}
		tx.Note, err = getField(models.TransactionNoteCol)
		if err != nil {
			return models.Transaction{}, err
		}
		return tx, nil
	}

	for {
		tx, err := getTransaction()
		if err != nil {
			return err
		}
		if _, err := db.InsertTransaction(tx); err != nil {
			return err
		}

		// TODO: Add context when adding transactions. e.g. making the last
		// used date the new default?, enabling an undo command
		// TODO: When adding a transaction, maybe show a couple of
		// transactions from around the same time?
		fmt.Print("\nWould you like to add another transaction? (y/[n]) ")
		confirmed, err := inpt.Confirm()
		fmt.Println()
		if err != nil {
			return err
		}
		if !confirmed {
			break
		}
	}

	return nil
}

func add(db *models.DB, cmdArgs []string) error {
	fs := flag.NewFlagSet(ingestName, flag.ContinueOnError)
	if err := fs.Parse(cmdArgs); err != nil {
		return err
	}
	args := fs.Args()
	if len(args) == 0 {
		return interactiveAdd(db)
	} else if len(args) > fieldsPerRecord {
		return fmt.Errorf("%s takes at most %d arguments", addName, fieldsPerRecord)
	}
	// TODO: implement an option that parses from flags or from args
	return nil
}
