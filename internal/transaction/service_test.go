package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kushturner/finances/internal/db"
	"github.com/stretchr/testify/assert"
)

type mockQuerier struct {
	transactions []db.Transaction
	err          error
}

func (m *mockQuerier) ListTransactions(ctx context.Context) ([]db.Transaction, error) {
	return m.transactions, m.err
}

func (m *mockQuerier) GetTransaction(ctx context.Context, id int32) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) UpdateTransaction(ctx context.Context, arg db.UpdateTransactionParams) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) DeleteTransaction(ctx context.Context, id int32) error {
	return nil
}

func (m *mockQuerier) CreateTransactionsBatch(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
	return 0, nil
}

func (m *mockQuerier) WithTx(tx interface{}) db.Querier {
	return m
}

func TestService_GetAllTransactions_EmptyList(t *testing.T) {
	mock := &mockQuerier{
		transactions: []db.Transaction{},
		err:          nil,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, transactions)
}

func TestService_GetAllTransactions_MultipleTransactions(t *testing.T) {
	category := "groceries"
	mock := &mockQuerier{
		transactions: []db.Transaction{
			{
				ID:          1,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "Grocery store",
				Amount:      5000,
				Currency:    "USD",
				Bank:        "Chase",
				Category:    pgtype.Text{String: category, Valid: true},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
			{
				ID:          2,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "Coffee shop",
				Amount:      500,
				Currency:    "GBP",
				Bank:        "Barclays",
				Category:    pgtype.Text{String: "", Valid: false},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		},
		err: nil,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	assert.Equal(t, int32(1), transactions[0].ID)
	assert.Equal(t, "Grocery store", transactions[0].Description)
	assert.Equal(t, int64(5000), transactions[0].Amount.Amount())
	assert.Equal(t, "USD", transactions[0].Amount.Currency().Code)
	assert.Equal(t, "Chase", transactions[0].Bank)
	assert.NotNil(t, transactions[0].Category)
	assert.Equal(t, "groceries", *transactions[0].Category)

	assert.Equal(t, int32(2), transactions[1].ID)
	assert.Equal(t, "Coffee shop", transactions[1].Description)
	assert.Equal(t, int64(500), transactions[1].Amount.Amount())
	assert.Equal(t, "GBP", transactions[1].Amount.Currency().Code)
	assert.Equal(t, "Barclays", transactions[1].Bank)
	assert.Nil(t, transactions[1].Category)
}

func TestService_GetAllTransactions_DatabaseError(t *testing.T) {
	mock := &mockQuerier{
		transactions: nil,
		err:          assert.AnError,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Equal(t, assert.AnError, err)
}

func TestService_GetAllTransactions_CategoryHandling(t *testing.T) {
	category := "transport"
	mock := &mockQuerier{
		transactions: []db.Transaction{
			{
				ID:          1,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "With category",
				Amount:      1000,
				Currency:    "USD",
				Bank:        "Test Bank",
				Category:    pgtype.Text{String: category, Valid: true},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
			{
				ID:          2,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "Without category",
				Amount:      2000,
				Currency:    "USD",
				Bank:        "Test Bank",
				Category:    pgtype.Text{String: "", Valid: false},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		},
		err: nil,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	assert.NotNil(t, transactions[0].Category)
	assert.Equal(t, "transport", *transactions[0].Category)

	assert.Nil(t, transactions[1].Category)
}
