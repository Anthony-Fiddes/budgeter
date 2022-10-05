package budgeter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/Anthony-Fiddes/budgeter/internal/csvfmt"
)

const (
	extCSV          = ".csv"
	fieldsPerRecord = 4
)

type ingest struct {
	Transactions Table
}

func newIngest(c *CLI) *ingest {
	return &ingest{Transactions: c.Transactions}
}

func (i ingest) Name() string {
	return "ingest"
}

//go:embed ingestUsage.txt
var ingestUsage string

func (i ingest) Usage() string {
	return ingestUsage
}

// ingest takes a file of valid transactions and inserts them all into the
// database.
//
// currently, it expects that the file type is included in the file name and
// only supports csv.
func (i ingest) Run(cmdArgs []string) error {
	// TODO: write tests
	// TODO: use a transaction so that all of the file is added or none of it is!
	fs := getFlagset(i.Name())
	err := fs.Parse(cmdArgs)
	if err != nil {
		return err
	}

	args := fs.Args()
	if len(args) != 1 {
		return fmt.Errorf("%s takes one argument", i.Name())
	}

	filePath := args[0]
	fileType := strings.ToLower(filepath.Ext(filePath))
	switch fileType {
	case extCSV:
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		b := csvfmt.NewBareReader(f)
		for {
			tx, err := b.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			err = i.Transactions.Insert(tx)
			if err != nil {
				return err
			}
		}
	case "":
		return fmt.Errorf("no file type specified")
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	return nil
}
