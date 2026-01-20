package transaction

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kushturner/finances/internal/db"
)

func TransactionFromDB(dbTx db.Transaction) Transaction {
	var category *string
	if dbTx.Category.Valid {
		category = &dbTx.Category.String
	}

	return Transaction{
		ID:          dbTx.ID,
		Date:        dbTx.Date.Time,
		Description: dbTx.Description,
		Amount:      dbTx.Amount,
		Currency:    dbTx.Currency,
		Bank:        dbTx.Bank,
		Category:    category,
		CreatedAt:   dbTx.CreatedAt.Time,
		UpdatedAt:   dbTx.UpdatedAt.Time,
	}
}

func TransactionToDB(tx Transaction) db.CreateTransactionParams {
	return db.CreateTransactionParams{
		Date:        pgtype.Date{Time: tx.Date, Valid: true},
		Description: tx.Description,
		Amount:      tx.Amount,
		Currency:    tx.Currency,
		Bank:        tx.Bank,
		Category:    pgtype.Text{String: stringOrEmpty(tx.Category), Valid: tx.Category != nil},
	}
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
