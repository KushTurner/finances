package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestListTransactions_EmptyList(t *testing.T) {
	mock := &mockQuerier{
		transactions: []db.Transaction{},
		err:          nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var responses []TransactionResponse
	err := json.NewDecoder(rec.Body).Decode(&responses)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(responses))
}

func TestListTransactions_MultipleTransactions(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var responses []TransactionResponse
	err := json.NewDecoder(rec.Body).Decode(&responses)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(responses))

	assert.Equal(t, int32(1), responses[0].ID)
	assert.Equal(t, "Grocery store", responses[0].Description)
	assert.Equal(t, "$50.00", responses[0].Amount)
	assert.Equal(t, "Chase", responses[0].Bank)
	assert.NotNil(t, responses[0].Category)
	assert.Equal(t, category, *responses[0].Category)

	assert.Equal(t, int32(2), responses[1].ID)
	assert.Equal(t, "£5.00", responses[1].Amount)
	assert.Nil(t, responses[1].Category)
}

func TestListTransactions_DatabaseError(t *testing.T) {
	mock := &mockQuerier{
		transactions: nil,
		err:          assert.AnError,
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListTransactions_CategoryNullable(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	var responses []TransactionResponse
	json.NewDecoder(rec.Body).Decode(&responses)

	assert.Equal(t, int32(1), responses[0].ID)
	assert.NotNil(t, responses[0].Category)
	assert.Equal(t, category, *responses[0].Category)
	assert.Equal(t, int32(2), responses[1].ID)
	assert.Nil(t, responses[1].Category)
}
