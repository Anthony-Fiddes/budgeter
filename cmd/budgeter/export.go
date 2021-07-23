package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const (
	exportName = "export"
)

// export writes all of the transactions in the given table to the given file name.
//
// currently, it expects that the file type is included in the file name and
// only supports csv.
// TODO: use a transaction so that all of the file is added or none of it is!
func export(table *transaction.Table, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(exportName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}

	args := fs.Args()
	if len(args) != 1 {
		return fmt.Errorf("%s takes one argument", exportName)
	}

	filePath := args[0]
	fileType := filepath.Ext(filePath)
	rows, err := table.Search("", -1)
	if err != nil {
		return err
	}
	switch fileType {
	case extCSV:
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		cw := transaction.NewCSVWriter(f)
		for rows.Next() {
			tx, err := rows.Scan()
			if err != nil {
				return err
			}
			cw.Write(tx)
		}
		cw.Flush()
	case "":
		return errors.New("no file type specified")
	default:
		return errors.New("unsupported file type")
	}

	fmt.Printf("Successfully exported to %s\n", filePath)
	return nil
}
