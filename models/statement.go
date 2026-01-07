package models

import (
	"time"
)

// Statement represents a bank statement for an account
type Statement struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	FromDate    time.Time `json:"from_date"`
	ToDate      time.Time `json:"to_date"`
	OpeningBalance float64   `json:"opening_balance"`
	ClosingBalance float64   `json:"closing_balance"`
	GeneratedAt time.Time `json:"generated_at"`
	Content     string    `json:"content"` // Could be a URL to the statement file or raw data
}
 

func (Statement) TableName() string {
	return "statements"
}

