package budgeter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

const exportName = "export"

// export writes all of the transactions in the given table to the given file name.
//
// currently, it expects that the file type is included in the file name and
// only supports csv.
func export(c *CLI) int {
	fs := getFlagset(exportName)
	err := fs.Parse(c.args)
	if err != nil {
		c.logParsingErr(err)
		return 1
	}

	args := fs.Args()
	if len(args) != 1 {
		c.err.Printf("%s takes one argument", exportName)
		return 1
	}

	filePath := args[0]
	fileType := strings.ToLower(filepath.Ext(filePath))
	rows, err := c.Transactions.Search("", -1)
	if err != nil {
		c.err.Println(err)
		return 1
	}
	switch fileType {
	case extCSV:
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			c.err.Println(err)
			return 1
		}
		defer f.Close()

		cw := transaction.NewCSVWriter(f)
		for rows.Next() {
			tx, err := rows.Scan()
			if err != nil {
				c.err.Println(err)
				return 1
			}
			cw.Write(tx)
		}
		cw.Flush()
	case "":
		c.err.Println("no file type specified")
		return 1
	default:
		c.err.Printf("unsupported file type: %s", fileType)
		return 1
	}

	return 0
}
