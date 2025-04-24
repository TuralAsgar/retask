package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Calculator *CalculatorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Calculator: &CalculatorModel{DB: db},
	}
}
