package transaction

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

const numCols = 4

type CSVWriter struct {
	*csv.Writer
}

// Write just writes a transaction to the internal buffer, you must remember to
// call Flush()
func (cw *CSVWriter) Write(tx Transaction) error {
	row := []string{
		tx.DateString(),
		tx.Entity,
		tx.Amount.String(),
		tx.Note,
	}
	return cw.Writer.Write(row)
}

func (cw *CSVWriter) WriteAll(txs []Transaction) error {
	for _, t := range txs {
		if err := cw.Write(t); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func NewCSVWriter(w io.Writer) *CSVWriter {
	cw := csv.NewWriter(w)
	return &CSVWriter{Writer: cw}
}

type CSVReader struct {
	*csv.Reader
}

// ? Should I consider allowing headers to set the order?
func (cr *CSVReader) Read() (Transaction, error) {
	cols, err := cr.Reader.Read()
	if err != nil {
		return Transaction{}, err
	}
	if len(cols) != numCols {
		row := strings.Join(cols, string(cr.Reader.Comma))
		return Transaction{}, fmt.Errorf(
			"transaction: CSV row \"%s\" must have %d columns",
			row, numCols,
		)
	}
	tx := Transaction{}
	tx.Date, err = Unix(cols[0])
	if err != nil {
		return Transaction{}, err
	}
	tx.Entity = cols[1]
	tx.Amount, err = GetCents(cols[2])
	if err != nil {
		return Transaction{}, err
	}
	tx.Note = cols[3]
	return tx, nil
}

func (cr *CSVReader) ReadAll() ([]Transaction, error) {
	var result []Transaction
	for {
		tx, err := cr.Read()
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

func NewCSVReader(r io.Reader) *CSVReader {
	cr := csv.NewReader(r)
	return &CSVReader{Reader: cr}
}
