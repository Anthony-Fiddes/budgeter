package budgeter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

type export struct {
	Transactions Table
}

func newExport(c *CLI) *export {
	return &export{Transactions: c.Transactions}
}

func (e export) Name() string {
	return "export"
}

func (e export) Usage() string {
	return "export writes all of your budgeter's transactions to a file. The file extension specified determines the format of the output."
}

// export writes all of the transactions in the given table to the given file name.
//
// currently, it expects that the file type is included in the file name and
// only supports csv.
func (e export) Run(cmdArgs []string) error {
	fs := getFlagset(e.Name())
	err := fs.Parse(cmdArgs)
	if err != nil {
		return err
	}

	args := fs.Args()
	if len(args) != 1 {
		return fmt.Errorf("%s takes one argument", e.Name())
	}

	filePath := args[0]
	fileType := strings.ToLower(filepath.Ext(filePath))
	rows, err := e.Transactions.Search("", -1)
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
		return fmt.Errorf("no file type specified")
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	return nil
}
