package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

		cw := transaction.NewCSVReader(f)
		for {
			tx, err := cw.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			err = table.Insert(tx)
			if err != nil {
				return err
			}
		}
	case "":
		return errors.New("no file type specified")
	default:
		return errors.New("unsupported file type")
	}

	fmt.Printf("Successfully ingested %s\n", filePath)
	return nil
}
