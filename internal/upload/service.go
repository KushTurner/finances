package upload

import (
	"context"
	"fmt"
	"io"

	"github.com/kushturner/finances/internal/csvparser"
	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/transaction"
)

type Service interface {
	UploadTransactions(ctx context.Context, bankType string, csvFile io.Reader) (int64, error)
}

type service struct {
	querier db.Querier
}

func NewService(querier db.Querier) Service {
	return &service{querier: querier}
}

func (s *service) UploadTransactions(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
	parser, err := csvparser.GetParser(bankType)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidBankType, err.Error())
	}

	transactions, err := parser.Parse(csvFile)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrParseFailure, err.Error())
	}

	batchParams := make([]db.CreateTransactionsBatchParams, len(transactions))
	for i, tx := range transactions {
		batchParams[i] = transaction.TransactionToBatchDB(tx)
	}

	count, err := s.querier.CreateTransactionsBatch(ctx, batchParams)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrDatabaseFailure, err.Error())
	}

	return count, nil
}
