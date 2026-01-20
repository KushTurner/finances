package transaction

import (
	"context"

	"github.com/kushturner/finances/internal/db"
)

type Service interface {
	GetAllTransactions(ctx context.Context) ([]Transaction, error)
}

type service struct {
	querier db.Querier
}

func NewService(querier db.Querier) Service {
	return &service{
		querier: querier,
	}
}

func (s *service) GetAllTransactions(ctx context.Context) ([]Transaction, error) {
	dbTransactions, err := s.querier.ListTransactions(ctx)
	if err != nil {
		return nil, err
	}

	transactions := make([]Transaction, 0, len(dbTransactions))
	for _, dbTx := range dbTransactions {
		transactions = append(transactions, TransactionFromDB(dbTx))
	}

	return transactions, nil
}
