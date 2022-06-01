package csvfmt

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Anthony-Fiddes/budgeter/model/transaction"
)

type Reader interface {
	Read() (transaction.Transaction, error)
}

func ReadAll(r Reader) ([]transaction.Transaction, error) {
	var result []transaction.Transaction
	for {
		tx, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		result = append(result, tx)
	}
	return result, nil
}

func getRow(cr *csv.Reader, fmtName string, numCols int) ([]string, error) {
	if cr == nil {
		panic("getRow cannot use a nil *csv.Reader")
	}
	cols, err := cr.Read()
	if err != nil {
		return nil, err
	}
	if len(cols) != numCols {
		row := strings.Join(cols, string(cr.Comma))
		return nil, fmt.Errorf(
			`transaction: CSV row "%s" must have %d columns in the %s format`,
			row, numCols, fmtName,
		)
	}
	return cols, nil
}

// TODO: Tests, tons of tests...

// Bare reads from a CSV file that expects transactions to be in the format of
// Date, Entity, Amount, Note.
type Bare struct {
	*csv.Reader
}

func NewBare(r io.Reader) *Bare {
	return &Bare{Reader: csv.NewReader(r)}
}

func (b *Bare) Read() (transaction.Transaction, error) {
	cols, err := getRow(b.Reader, "bare", 4)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx := transaction.Transaction{}
	tx.Date, err = transaction.Unix(cols[0])
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Entity = cols[1]
	tx.Amount, err = transaction.GetCents(cols[2])
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Note = cols[3]
	return tx, nil
}

const Chase = "chase"

type ChaseReader struct {
	*csv.Reader
}

func NewChaseReader(r io.Reader) *ChaseReader {
	c := &ChaseReader{Reader: csv.NewReader(r)}
	// discard header
	c.Reader.Read()
	return c
}

func (c *ChaseReader) Read() (transaction.Transaction, error) {
	cols, err := getRow(c.Reader, Chase, 7)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx := transaction.Transaction{}
	tx.Date, err = transaction.Unix(cols[1])
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Entity = cols[2]
	tx.Amount, err = transaction.GetCents(cols[5])
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Note = cols[6]
	return tx, nil
}

func convDate(dateStr string, currLayout, newLayout string) (string, error) {
	date, err := time.Parse(currLayout, dateStr)
	if err != nil {
		return "", err
	}
	return date.Format(newLayout), nil
}

const CapitalOne = "capitalone"

type CapitalOneReader struct {
	*csv.Reader
}

func NewCapitalOneReader(r io.Reader) *CapitalOneReader {
	co := &CapitalOneReader{Reader: csv.NewReader(r)}
	// discard header line
	co.Reader.Read()
	return co
}

func (c *CapitalOneReader) Read() (transaction.Transaction, error) {
	cols, err := getRow(c.Reader, CapitalOne, 7)
	if err != nil {
		return transaction.Transaction{}, err
	}
	const dateLayout = "2006-01-02"
	date, err := convDate(cols[1], dateLayout, transaction.DateLayout)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx := transaction.Transaction{}
	tx.Date, err = transaction.Unix(date)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx.Entity = cols[3]
	if cols[5] != "" {
		tx.Amount, err = transaction.GetCents(cols[5])
		if err != nil {
			return transaction.Transaction{}, err
		}
		tx.Amount *= -1
	} else {
		tx.Amount, err = transaction.GetCents(cols[6])
	}
	return tx, nil
}
