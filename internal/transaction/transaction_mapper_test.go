package transaction

import (
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kushturner/finances/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestTransactionFromDB_WithValidCategory(t *testing.T) {
	categoryStr := "groceries"
	dbTx := db.Transaction{
		ID:          1,
		Date:        pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
		Description: "Test transaction",
		Amount:      5000,
		Currency:    "USD",
		Bank:        "Chase",
		Category:    pgtype.Text{String: categoryStr, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	result := TransactionFromDB(dbTx)

	assert.NotNil(t, result.Category)
	assert.Equal(t, categoryStr, *result.Category)
	assert.Equal(t, dbTx.ID, result.ID)
	assert.Equal(t, dbTx.Amount, result.Amount.Amount())
	assert.Equal(t, dbTx.Currency, result.Amount.Currency().Code)
}

func TestTransactionFromDB_WithNullCategory(t *testing.T) {
	dbTx := db.Transaction{
		ID:          1,
		Date:        pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
		Description: "Test transaction",
		Amount:      5000,
		Currency:    "USD",
		Bank:        "Chase",
		Category:    pgtype.Text{String: "", Valid: false},
		CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	result := TransactionFromDB(dbTx)

	assert.Nil(t, result.Category)
}

func TestTransactionToDB_WithCategory(t *testing.T) {
	categoryStr := "groceries"
	tx := Transaction{
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Description: "Test transaction",
		Amount:      money.New(5000, "USD"),
		Bank:        "Chase",
		Category:    &categoryStr,
	}

	result := TransactionToDB(tx)

	assert.True(t, result.Category.Valid)
	assert.Equal(t, categoryStr, result.Category.String)
	assert.Equal(t, int64(5000), result.Amount)
	assert.Equal(t, "USD", result.Currency)
}

func TestTransactionToDB_WithNilCategory(t *testing.T) {
	tx := Transaction{
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Description: "Test transaction",
		Amount:      money.New(5000, "USD"),
		Bank:        "Chase",
		Category:    nil,
	}

	result := TransactionToDB(tx)

	assert.False(t, result.Category.Valid)
	assert.Equal(t, "", result.Category.String)
}

func TestStringOrEmpty_WithNil(t *testing.T) {
	result := stringOrEmpty(nil)

	assert.Equal(t, "", result)
}

func TestStringOrEmpty_WithValue(t *testing.T) {
	str := "test value"
	result := stringOrEmpty(&str)

	assert.Equal(t, str, result)
}
