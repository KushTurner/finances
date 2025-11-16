package transactions

import (
	"context"
	"testing"

	repo "github.com/kushturner/finances/internal/database/psql/sqlc"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	txs []repo.Transaction
	err error
}

func (m *mockRepository) GetAllTransactions(ctx context.Context) ([]repo.Transaction, error) {
	return m.txs, m.err
}

func TestListTransactions(t *testing.T) {
	t.Run("returns all transactions from repository", func(t *testing.T) {
		expected := []repo.Transaction{
			{Description: "tx 1"},
			{Description: "tx 2"},
		}
		repo := &mockRepository{txs: expected}
		svc := NewService(repo)

		actual, _ := svc.ListTransactions(context.Background())

		assert.Len(t, actual, 2)
	})

	t.Run("returns empty slice if no transactions found in repository", func(t *testing.T) {
		repo := &mockRepository{}
		svc := NewService(repo)

		actual, _ := svc.ListTransactions(context.Background())

		assert.Len(t, actual, 0)
	})
}
