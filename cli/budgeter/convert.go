package budgeter

import (
	"fmt"
	"os"
	"path"

	"github.com/Anthony-Fiddes/budgeter/internal/csvfmt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

type convert struct{}

func (c convert) Name() string {
	return "convert"
}

func (c convert) Usage() string {
	return "convert takes a CSV from a financial provider and simplifies it to the budgeter format"
}

func (c convert) Run(cmdArgs []string) error {
	fs := getFlagset(c.Name())
	if err := fs.Parse(cmdArgs); err != nil {
		return err
	}
	// TODO: guess the file format from its name or other details
	args := fs.Args()
	if len(args) < 2 {
		return fmt.Errorf("%s takes at least two arguments (the name of the file to be converted and its format)", c.Name())
	}

	inputPath := args[0]
	format := args[1]
	f, err := os.Open(args[0])
	defer f.Close()
	if err != nil {
		return err
	}

	var transactions []transaction.Transaction
	switch format {
	case csvfmt.Chase:
		chase := csvfmt.NewChaseReader(f)
		transactions, err = csvfmt.ReadAll(chase)
		if err != nil {
			return fmt.Errorf(`could not read transactions from %s formatted file %s: %w`, csvfmt.Chase, inputPath, err)
		}
	case csvfmt.CapitalOne:
		capitalone := csvfmt.NewCapitalOneReader(f)
		transactions, err = csvfmt.ReadAll(capitalone)
		if err != nil {
			return fmt.Errorf(`could not read transactions from %s formatted file %s: %w`, csvfmt.CapitalOne, inputPath, err)
		}
	case csvfmt.Venmo:
		venmo := csvfmt.NewVenmoReader(f)
		transactions, err = csvfmt.ReadAll(venmo)
		if err != nil {
			return fmt.Errorf(`could not read transactions from %s formatted file %s: %w`, csvfmt.Venmo, inputPath, err)
		}
	default:
		// TODO: Consider listing supported formats
		return fmt.Errorf(`format "%s" is not supported`, format)
	}

	outputPath := path.Join(path.Dir(inputPath), "CONVERTED_"+path.Base(inputPath))
	out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer out.Close()
	w := transaction.NewCSVWriter(out)
	err = w.WriteAll(transactions)
	if err != nil {
		return fmt.Errorf(`could not write converted file to %s: %v`, outputPath, err)
	}
	return nil
}
