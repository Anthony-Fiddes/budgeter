package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
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
func ingest(db *models.DB, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(ingestName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}

	args := fs.Args()
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}

	filePath := args[0]
	fileType := filepath.Ext(filePath)
	switch fileType {
	case extCSV:
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = ingestCSV(f, db)
		if err != nil {
			return err
		}
	case "":
		return errors.New("no file type specified")
	default:
		return errors.New("unsupported file type")
	}

	return nil
}

func ingestCSV(r io.Reader, db *models.DB) error {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = fieldsPerRecord
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		tx, err := csvRowToTx(row)
		if err != nil {
			return err
		}
		_, err = db.InsertTransaction(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func csvRowToTx(row []string) (models.Transaction, error) {
	if len(row) < 4 {
		return models.Transaction{}, errors.New(`not enough columns in input`)
	}
	for i := range row {
		row[i] = strings.TrimSpace(row[i])
	}

	// TODO: figure out whether or not I want to handle weirdly formatted
	// amounts. e.g. $5 instead of $5.00
	amount := row[2]
	amount = strings.Replace(amount, ".", "", 1)
	amount = strings.Replace(amount, string(models.Currency), "", 1)
	amount = strings.Replace(amount, ",", "", 1)
	a, err := strconv.Atoi(amount)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error parsing the amount for a transaction: %w", err)
	}
	date := row[0]
	d, err := time.Parse("1/2/2006", date)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error parsing the date for a transaction: %w", err)
	}
	tx := models.Transaction{
		Entity: row[1],
		Amount: a,
		Date:   d.Unix(),
		Note:   row[3],
	}
	return tx, err
}
