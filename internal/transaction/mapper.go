package transaction

import (
	"github.com/Rhymond/go-money"
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
		Amount:      money.New(dbTx.Amount, dbTx.Currency),
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
		Amount:      tx.Amount.Amount(),
		Currency:    tx.Amount.Currency().Code,
		Bank:        tx.Bank,
		Category:    pgtype.Text{String: stringOrEmpty(tx.Category), Valid: tx.Category != nil},
	}
}

func TransactionToBatchDB(tx Transaction) db.CreateTransactionsBatchParams {
	return db.CreateTransactionsBatchParams{
		Date:        pgtype.Date{Time: tx.Date, Valid: true},
		Description: tx.Description,
		Amount:      tx.Amount.Amount(),
		Currency:    tx.Amount.Currency().Code,
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
