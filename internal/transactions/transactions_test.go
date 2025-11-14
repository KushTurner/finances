package transactions

import "testing"

type mockRepository struct {
	txs []Transaction
	err error
}

func (m *mockRepository) GetAllTransactions() ([]Transaction, error) {
	return m.txs, m.err
}

func TestListTransactions(t *testing.T) {
	t.Run("returns all transactions from repository", func(t *testing.T) {
		expected := []Transaction{
			{Description: "tx 1"},
			{Description: "tx 2"},
		}
		repo := &mockRepository{txs: expected}
		svc := NewService(repo)
		actual, _ := svc.ListTransactions()

		for i := range expected {
			if actual[i].Description != expected[i].Description {
				t.Errorf("got %v, want %v", actual[i].Description, expected[i].Description)
			}
		}
	})

	t.Run("returns empty slice if no transactions found in repository", func(t *testing.T) {
		repo := &mockRepository{}
		svc := NewService(repo)
		actual, _ := svc.ListTransactions()

		if len(actual) != 0 {
			t.Errorf("got %d transactions, want 0", len(actual))
		}
	})
}
