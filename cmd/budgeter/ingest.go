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
func ingest(db *models.DB, cmdArgs []string) error {
	var err error
	fs := flag.NewFlagSet(ingestName, flag.ContinueOnError)
	err = fs.Parse(cmdArgs)
	if err != nil {
		return err
	}
	args := fs.Args()

	filePath := args[0]
	fileType := filepath.Ext(filePath)
	switch fileType {
	case extCSV:
		err = ingestCSV(filePath, db)
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

func ingestCSV(filePath string, db *models.DB) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	cr := csv.NewReader(f)
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

	amount := row[2]
	amount = strings.Replace(amount, ".", "", 1)
	amount = strings.Replace(amount, "$", "", 1)
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