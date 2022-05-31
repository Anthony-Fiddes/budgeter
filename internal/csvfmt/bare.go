package csvfmt

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

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

// Bare reads from a CSV file that expects transactions to be in the format of
// Date, Entity, Amount, Note.
type Bare struct {
	*csv.Reader
}

func NewBare(r io.Reader) *Bare {
	return &Bare{Reader: csv.NewReader(r)}
}

func (b *Bare) Read() (transaction.Transaction, error) {
	const numCols = 4
	cols, err := b.Reader.Read()
	if err != nil {
		return transaction.Transaction{}, err
	}
	if len(cols) != numCols {
		row := strings.Join(cols, string(b.Reader.Comma))
		return transaction.Transaction{}, fmt.Errorf(
			"transaction: CSV row \"%s\" must have %d columns",
			row, numCols,
		)
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

type Chase struct {
	*csv.Reader
}

func NewChase(r io.Reader) *Chase {
	return &Chase{Reader: csv.NewReader(r)}
}

func (c *Chase) Read() (transaction.Transaction, error) {
	const fmtName = "chase"
	const numCols = 7
	c.Reader.Read()
	cols, err := c.Reader.Read()
	if err != nil {
		return transaction.Transaction{}, err
	}
	if len(cols) != numCols {
		row := strings.Join(cols, string(c.Reader.Comma))
		return transaction.Transaction{}, fmt.Errorf(
			"transaction: CSV row \"%s\" must have %d columns in the %s format",
			row, numCols, fmtName,
		)
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
