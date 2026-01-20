package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/kushturner/finances/internal/transaction"
	"github.com/stretchr/testify/assert"
)

type mockTransactionService struct {
	transactions []transaction.Transaction
	err          error
}

func (m *mockTransactionService) GetAllTransactions(ctx context.Context) ([]transaction.Transaction, error) {
	return m.transactions, m.err
}

func TestListTransactions_EmptyList(t *testing.T) {
	mock := &mockTransactionService{
		transactions: []transaction.Transaction{},
		err:          nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, "[]", rec.Body.String())
}

func TestListTransactions_MultipleTransactions(t *testing.T) {
	category := "groceries"
	mock := &mockTransactionService{
		transactions: []transaction.Transaction{
			{
				ID:          1,
				Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Description: "Grocery store",
				Amount:      money.New(5000, "USD"),
				Bank:        "Chase",
				Category:    &category,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				Date:        time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				Description: "Coffee shop",
				Amount:      money.New(500, "GBP"),
				Bank:        "Barclays",
				Category:    nil,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		err: nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedJSON := `[
		{
			"id": 1,
			"date": "2024-01-15T00:00:00Z",
			"description": "Grocery store",
			"amount": "$50.00",
			"bank": "Chase",
			"category": "groceries"
		},
		{
			"id": 2,
			"date": "2024-01-16T00:00:00Z",
			"description": "Coffee shop",
			"amount": "£5.00",
			"bank": "Barclays",
			"category": null
		}
	]`

	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

func TestListTransactions_DatabaseError(t *testing.T) {
	mock := &mockTransactionService{
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
	mock := &mockTransactionService{
		transactions: []transaction.Transaction{
			{
				ID:          1,
				Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Description: "With category",
				Amount:      money.New(1000, "USD"),
				Bank:        "Test Bank",
				Category:    &category,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				Date:        time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				Description: "Without category",
				Amount:      money.New(2000, "USD"),
				Bank:        "Test Bank",
				Category:    nil,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		err: nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	rec := httptest.NewRecorder()

	handler := NewListTransactionsHandler(mock)
	handler(rec, req)

	expectedJSON := `[
		{
			"id": 1,
			"date": "2024-01-15T00:00:00Z",
			"description": "With category",
			"amount": "$10.00",
			"bank": "Test Bank",
			"category": "transport"
		},
		{
			"id": 2,
			"date": "2024-01-16T00:00:00Z",
			"description": "Without category",
			"amount": "$20.00",
			"bank": "Test Bank",
			"category": null
		}
	]`

	assert.JSONEq(t, expectedJSON, rec.Body.String())
}
