package budgeter

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/internal/csvfmt"
)

const (
	ingestName = "ingest"
	extCSV     = ".csv"
)

//go:embed ingestUsage.txt
var ingestUsage string

// ingest takes a file of valid transactions and inserts them all into the
// database.
//
// currently, it expects that the file type is included in the file name and
// only supports csv.
func ingest(c *CLI) int {
	// TODO: write tests
	// TODO: use a transaction so that all of the file is added or none of it is!
	fs := getFlagset(ingestName)
	err := fs.Parse(c.args)
	if err != nil {
		c.logParsingErr(err)
		return 1
	}

	args := fs.Args()
	if len(args) != 1 {
		c.err.Printf("%s takes one argument", ingestName)
		c.err.Println()
		c.err.Println(ingestUsage)
		return 1
	}

	filePath := args[0]
	fileType := strings.ToLower(filepath.Ext(filePath))
	switch fileType {
	case extCSV:
		f, err := os.Open(filePath)
		if err != nil {
			c.err.Printf("could not open \"%s\": %v", filePath, err)
			return 1
		}
		defer f.Close()

		b := csvfmt.NewBare(f)
		for {
			tx, err := b.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				c.err.Println(err)
				return 1
			}
			err = c.Transactions.Insert(tx)
			if err != nil {
				c.err.Println(err)
				return 1
			}
		}
	case "":
		c.err.Println("no file type specified")
		return 1
	default:
		c.err.Printf("unsupported file type: %s", fileType)
		c.err.Println()
		c.err.Println(ingestUsage)
		return 1
	}

	return 0
}
