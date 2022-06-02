package budgeter

import (
	"os"
	"path"

	"github.com/Anthony-Fiddes/budgeter/internal/csvfmt"
	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

var convertUsage = "convert takes a CSV from a financial provider and simplifies it to the budgeter format"

type convert struct{}

func (c *convert) Name() string {
	return "convert"
}

func (c *convert) Run(cli *CLI) int {
	fs := getFlagset(c.Name())
	if err := fs.Parse(cli.args); err != nil {
		cli.logParsingErr(err)
		cli.err.Println()
		cli.err.Print(convertUsage)
		return 1
	}
	// TODO: guess the file format from its name or other details
	args := fs.Args()
	if len(args) < 2 {
		cli.err.Printf("%s takes at least two arguments (the name of the file to be converted and its format)", c.Name())
		return 1
	}

	inputPath := args[0]
	format := args[1]
	f, err := os.Open(args[0])
	defer f.Close()
	if err != nil {
		cli.err.Println(err)
		return 1
	}

	var transactions []transaction.Transaction
	switch format {
	case csvfmt.Chase:
		chase := csvfmt.NewChaseReader(f)
		transactions, err = csvfmt.ReadAll(chase)
		if err != nil {
			cli.err.Printf(`could not read transactions from %s formatted file %s: %v`, csvfmt.Chase, inputPath, err)
		}
	case csvfmt.CapitalOne:
		capitalone := csvfmt.NewCapitalOneReader(f)
		transactions, err = csvfmt.ReadAll(capitalone)
		if err != nil {
			cli.err.Printf(`could not read transactions from %s formatted file %s: %v`, csvfmt.CapitalOne, inputPath, err)
		}
	case csvfmt.Venmo:
		venmo := csvfmt.NewVenmoReader(f)
		transactions, err = csvfmt.ReadAll(venmo)
		if err != nil {
			cli.err.Printf(`could not read transactions from %s formatted file %s: %v`, csvfmt.Venmo, inputPath, err)
		}
	default:
		// TODO: List supported formats
		cli.err.Printf(`format "%s" is not supported`, format)
	}

	outputPath := path.Join(path.Dir(inputPath), "CONVERTED_"+path.Base(inputPath))
	out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer out.Close()
	w := transaction.NewCSVWriter(out)
	err = w.WriteAll(transactions)
	if err != nil {
		cli.err.Printf(`could not write converted file to %s: %v`, outputPath, err)
	}
	return 0
}
