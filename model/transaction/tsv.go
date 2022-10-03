package transaction

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type TSVTable struct {
	filePath string
}

func NewTSVTable(file string) *TSVTable {
	t := &TSVTable{filePath: file}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(file, os.O_CREATE, 0644)
		f.Close()
		if err != nil {
			log.Panicf(`Could not create TSVTable backing file at %v: %v`, file, err)
		}
	}
	return t
}

func (t TSVTable) Insert(tx Transaction) error {
	file, err := os.OpenFile(t.filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("could not open TSV file for append: %w", err)
	}
	defer file.Close()
	tw := newTSVWriter(file)
	err = tw.Write(tx)
	if err != nil {
		return fmt.Errorf("TSVTable.Insert: could not add transaction to TSV file: %w", err)
	}
	return nil
}

func (t TSVTable) RangeTotal(start, end time.Time) (Cent, error) {
	transactions, err := t.read()
	if err != nil {
		return 0, fmt.Errorf("TSVTable.RangeTotal: %w", err)
	}
	startUnix := start.Unix()
	endUnix := end.Unix()
	var sum Cent
	for _, tx := range transactions {
		if tx.Date >= startUnix && tx.Date <= endUnix {
			sum += tx.Amount
		}
	}
	return sum, nil
}

func (t TSVTable) Remove(transactionID int) error {
	transactions, err := t.read()
	if err != nil {
		return fmt.Errorf("TSVTable.Remove: %w", err)
	}
	for i, tx := range transactions {
		if tx.ID == transactionID {
			transactions[i] = transactions[len(transactions)-1]
			transactions = transactions[:len(transactions)-1]
			err := t.write(transactions)
			if err != nil {
				return fmt.Errorf("TSVTable.Remove: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (t TSVTable) Search(query string, limit int) ([]Transaction, error) {
	transactions, err := t.read()
	if err != nil {
		return nil, fmt.Errorf("TSVTable.Search: %w", err)
	}
	var result []Transaction
	query = strings.ToLower(query)
	for _, tx := range transactions {
		if limit > 0 && len(result) >= limit {
			break
		}
		if strings.Contains(strings.ToLower(tx.Entity), query) || strings.Contains(strings.ToLower(tx.Note), query) {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (t TSVTable) read() ([]Transaction, error) {
	file, err := os.Open(t.filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open TSV file for read: %w", err)
	}
	defer file.Close()
	tr := newTSVReader(file)
	transactions, err := tr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read all transactions from TSV file: %w", err)
	}
	for i := range transactions {
		transactions[i].ID = i
	}
	return transactions, nil
}

func (t TSVTable) write(transactions []Transaction) error {
	file, err := os.OpenFile(t.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open TSV file for write: %w", err)
	}
	defer file.Close()
	tw := newTSVWriter(file)
	err = tw.WriteAll(transactions)
	if err != nil {
		return fmt.Errorf("could not write all transactions to TSV file: %w", err)
	}
	return nil
}

func newTSVReader(r io.Reader) *CSVReader {
	csvReader := NewCSVReader(r)
	csvReader.Comma = '\t'
	csvReader.Read()
	return csvReader
}

func newTSVWriter(w io.Writer) *CSVWriter {
	csvWriter := NewCSVWriter(w)
	csvWriter.Comma = '\t'
	return csvWriter
}
