package main

import (
	"flag"
	"fmt"

	"github.com/Anthony-Fiddes/budgeter/internal/inpt"
	"github.com/Anthony-Fiddes/budgeter/internal/models"
)

const addName = "add"

// TODO: Find a way to handle duplicates gracefully
func interactiveAdd() (models.Transaction, error) {
	addField := func(field string) (string, error) {
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
	tx := models.Transaction{}
	date, err := addField(models.TransactionDateCol)
	if err != nil {
		return models.Transaction{}, err
	}
	tx.Date, err = models.Date(date)
	if err != nil {
		return models.Transaction{}, err
	}
	tx.Entity, err = addField(models.TransactionEntityCol)
	if err != nil {
		return models.Transaction{}, err
	}
	amount, err := addField(models.TransactionAmountCol)
	if err != nil {
		return models.Transaction{}, err
	}
	tx.Amount, err = models.Cents(amount)
	if err != nil {
		return models.Transaction{}, err
	}
	tx.Note, err = addField(models.TransactionNoteCol)
	if err != nil {
		return models.Transaction{}, err
	}
	return tx, nil
}

func add(db *models.DB, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(ingestName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}
	args := fs.Args()
	if len(args) == 0 {
		for {
			tx, err := interactiveAdd()
			if err != nil {
				return err
			}
			if _, err := db.InsertTransaction(tx); err != nil {
				return err
			}

			// TODO: Add context when adding transactions. e.g. making the last
			// used date the new default?, enabling an undo command
			fmt.Print("\nWould you like to add another transaction? ")
			confirmed, err := inpt.Confirm()
			fmt.Println()
			if err != nil {
				return err
			}
			if !confirmed {
				break
			}
		}
	}
	// TODO: implement an option that parses from flags or from args
	return nil
}
