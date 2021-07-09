package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const (
	ingestName      = "ingest"
	extCSV          = ".csv"
	fieldsPerRecord = 4
)

// ingest takes a file of valid transactions and inserts them all into the
// database.
//
// currently, it expects that the file type is included in the file name and
// only supports csv.
// TODO: write tests
// TODO: use a transaction so that all of the file is added or none of it is!
func ingest(table *transaction.Table, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(ingestName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}

	args := fs.Args()
	if len(args) != 1 {
		return fmt.Errorf("%s takes one argument", ingestName)
	}

	filePath := args[0]
	fileType := filepath.Ext(filePath)
	switch fileType {
	case extCSV:
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = ingestCSV(f, table)
		if err != nil {
			return err
		}
	case "":
		return errors.New("no file type specified")
	default:
		return errors.New("unsupported file type")
	}

	fmt.Printf("Successfully ingested %s\n", filePath)
	return nil
}

// TODO: Inserter and ingestCSV should be refactored out to a package with tests
type Inserter interface {
	Insert(transaction.Transaction) error
}

func ingestCSV(r io.Reader, in Inserter) error {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = fieldsPerRecord
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		tx, err := csvRowToTx(row)
		if err != nil {
			return err
		}
		err = in.Insert(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func csvRowToTx(row []string) (transaction.Transaction, error) {
	if len(row) < 4 {
		return transaction.Transaction{}, errors.New(`not enough columns in input`)
	}
	for i := range row {
		row[i] = strings.TrimSpace(row[i])
	}

	amount, err := transaction.Cents(row[2])
	if err != nil {
		return transaction.Transaction{}, err
	}
	d := row[0]
	date, err := transaction.Date(d)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx := transaction.Transaction{
		Entity: row[1],
		Amount: amount,
		Date:   date,
		Note:   row[3],
	}
	return tx, err
}
