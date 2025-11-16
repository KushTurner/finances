package transactions

import (
	"context"
	"time"

	"github.com/Rhymond/go-money"
	repo "github.com/kushturner/finances/internal/database/psql/sqlc"
)

type Transaction struct {
	Date        time.Time
	Description string
	Out         *money.Money
	In          *money.Money
}

type Service interface {
	ListTransactions(ctx context.Context) ([]Transaction, error)
}

type repository interface {
	GetAllTransactions(ctx context.Context) ([]repo.Transaction, error)
}

type service struct {
	repo repository
}

func (s *service) ListTransactions(ctx context.Context) ([]Transaction, error) {
	txs, err := s.repo.GetAllTransactions(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]Transaction, len(txs))
	for i, tx := range txs {
		result[i] = toTransaction(tx)
	}

	return result, nil
}

func NewService(repo repository) Service {
	return &service{repo: repo}
}

func toTransaction(tx repo.Transaction) Transaction {
	return Transaction{
		Date:        tx.Date.Time,
		Description: tx.Description,
		Out:         money.New(tx.AmountOut.Int64, tx.Currency),
		In:          money.New(tx.AmountIn.Int64, tx.Currency),
	}
}
