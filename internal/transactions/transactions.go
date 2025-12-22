package transactions

import (
	"context"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/kushturner/finances/internal/database/psql/sqlc"
)

type Transaction struct {
	Date        time.Time
	Description string
	Out         *money.Money
	In          *money.Money
	Bank        string
}

type Service interface {
	ListTransactions(ctx context.Context) ([]Transaction, error)
	AddTransactions(ctx context.Context, txs []Transaction) ([]Transaction, error)
}

type repository interface {
	GetAllTransactions(ctx context.Context) ([]repo.Transaction, error)
	AddTransaction(ctx context.Context, arg repo.AddTransactionParams) (repo.Transaction, error)
}

type service struct {
	repo repository
}

func NewService(repo repository) Service {
	return &service{repo: repo}
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

func (s *service) AddTransactions(ctx context.Context, txs []Transaction) ([]Transaction, error) {
	var result []Transaction
	for _, tx := range txs {

		arg := repo.AddTransactionParams{
			Date:        pgtype.Date{Time: tx.Date, Valid: true},
			Description: tx.Description,
			AmountOut:   pgtype.Int8{Int64: tx.Out.Amount(), Valid: true},
			AmountIn:    pgtype.Int8{Int64: tx.In.Amount(), Valid: true},
			Currency:    tx.Out.Currency().Code,
			Bank:        pgtype.Text{String: tx.Bank, Valid: true},
		}

		_, err := s.repo.AddTransaction(ctx, arg)
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}

	return result, nil
}

func toTransaction(tx repo.Transaction) Transaction {
	return Transaction{
		Date:        tx.Date.Time,
		Description: tx.Description,
		Out:         money.New(tx.AmountOut.Int64, tx.Currency),
		In:          money.New(tx.AmountIn.Int64, tx.Currency),
		Bank:        tx.Bank.String,
	}
}
