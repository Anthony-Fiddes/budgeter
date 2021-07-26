package budgeter

import (
	"flag"
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
func ingest(c *CLI) int {
	fs := flag.NewFlagSet(ingestName, flag.ContinueOnError)
	err := fs.Parse(c.args)
	if err != nil {
		c.logParsingErr(err)
		return 1
	}

	args := fs.Args()
	if len(args) != 1 {
		c.Log.Printf("%s takes one argument", ingestName)
		return 1
	}

	filePath := args[0]
	fileType := filepath.Ext(filePath)
	switch fileType {
	case extCSV:
		f, err := os.Open(filePath)
		if err != nil {
			c.Log.Printf("could not open \"%s\": %v", filePath, err)
			return 1
		}
		defer f.Close()

		cw := transaction.NewCSVReader(f)
		for {
			tx, err := cw.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				c.Log.Println(err)
				return 1
			}
			err = c.Transactions.Insert(tx)
			if err != nil {
				c.Log.Println(err)
				return 1
			}
		}
	case "":
		c.Log.Println("no file type specified")
		return 1
	default:
		c.Log.Println("unsupported file type")
		return 1
	}

	return 0
}
