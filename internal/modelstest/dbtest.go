package modelstest

import (
	"database/sql"
	"fmt"

	"github.com/Anthony-Fiddes/budgeter/internal/models"
)

const (
	sqlite3URI = ":memory:"
)

func GetMemDB() (*models.DB, error) {
	db, err := sql.Open("sqlite3", sqlite3URI)
	if err != nil {
		return nil, fmt.Errorf("error creating an in-memory database for testing: %w", err)
	}
	return &models.DB{DB: db}, nil
}

func GetMemDBWithTable() (*models.DB, error) {
	db, err := GetMemDB()
	if err != nil {
		return nil, err
	}
	_, err = db.CreateTransactionTable()
	if err != nil {
		return nil, fmt.Errorf("error creating the transaction table: %w", err)
	}
	return db, nil
}
