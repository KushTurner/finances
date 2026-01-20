package handlers

import (
	"time"

	"github.com/kushturner/finances/internal/transaction"
)

type TransactionResponse struct {
	ID          int32     `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      string    `json:"amount"`
	Bank        string    `json:"bank"`
	Category    *string   `json:"category"`
}

func FromTransaction(t transaction.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          t.ID,
		Date:        t.Date,
		Description: t.Description,
		Amount:      t.Amount.Display(),
		Bank:        t.Bank,
		Category:    t.Category,
	}
}
