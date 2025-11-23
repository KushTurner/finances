package handlers

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/kushturner/finances/internal/transactions"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	txs []transactions.Transaction
	err error
}

func (m *mockService) ListTransactions(ctx context.Context) ([]transactions.Transaction, error) {
	return m.txs, m.err
}

func (m *mockService) AddTransactions(ctx context.Context, txs []transactions.Transaction) ([]transactions.Transaction, error) {
	m.txs = append(m.txs, txs...)
	return txs, m.err
}

func TestListTransactions(t *testing.T) {
	t.Run("can get all transactions", func(t *testing.T) {
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		mockTx := transactions.Transaction{
			Date:        testDate,
			Description: "Chrome Hearts Ring",
			Out:         money.New(26700, "GBP"),
			In:          money.New(0, "GBP"),
		}
		svc := &mockService{txs: []transactions.Transaction{mockTx}}
		handler := NewTransactionHandler(svc)

		expected := `[{"date":"2024-01-15T00:00:00Z","description":"Chrome Hearts Ring","out":"£267.00","in":"£0.00"}]`

		actual := assert.HTTPBody(handler.ListTransactions, http.MethodGet, "/transactions", nil)

		assert.JSONEq(t, expected, actual)
	})

	t.Run("returns empty list if no transactions exist", func(t *testing.T) {
		svc := &mockService{}
		handler := NewTransactionHandler(svc)

		expected := `[]`

		actual := assert.HTTPBody(handler.ListTransactions, http.MethodGet, "/transactions", nil)

		assert.JSONEq(t, expected, actual)
	})
}
