package transaction

import (
	"context"
	"fmt"
	"io"

	"github.com/kushturner/finances/internal/db"
)

type Service interface {
	GetAllTransactions(ctx context.Context) ([]Transaction, error)
	ImportFromCSV(ctx context.Context, parser Parser, csvFile io.Reader) (int64, error)
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

func (s *service) ImportFromCSV(ctx context.Context, parser Parser, csvFile io.Reader) (int64, error) {
	transactions, err := parser.Parse(csvFile)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrParseFailure, err.Error())
	}

	batchParams := make([]db.CreateTransactionsBatchParams, len(transactions))
	for i, tx := range transactions {
		batchParams[i] = TransactionToBatchDB(tx)
	}

	count, err := s.querier.CreateTransactionsBatch(ctx, batchParams)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrDatabaseFailure, err.Error())
	}

	return count, nil
}
