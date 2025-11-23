package transactions

import (
	"context"
	"testing"

	"github.com/Rhymond/go-money"
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

func (m *mockRepository) AddTransaction(ctx context.Context, arg repo.AddTransactionParams) (repo.Transaction, error) {
	m.txs = append(m.txs, repo.Transaction{})
	return repo.Transaction{}, nil
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

func TestAddTransaction(t *testing.T) {
	t.Run("adds multiple transactions to the repository", func(t *testing.T) {
		repo := &mockRepository{}
		svc := NewService(repo)

		tx := Transaction{Description: "new tx", Out: money.New(1000, "GBP"), In: money.New(0, "GBP")}
		tx2 := Transaction{Description: "new tx 2", Out: money.New(2000, "GBP"), In: money.New(0, "GBP")}

		svc.AddTransactions(context.Background(), []Transaction{tx, tx2})

		assert.Len(t, repo.txs, 2)
	})
}
