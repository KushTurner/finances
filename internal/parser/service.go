package parser

import (
	"fmt"
	"io"

	"github.com/kushturner/finances/internal/csvparser"
	"github.com/kushturner/finances/internal/transaction"
)

type Service interface {
	Parse(r io.Reader, bankType string) ([]transaction.Transaction, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Parse(r io.Reader, bankType string) ([]transaction.Transaction, error) {
	parser, err := csvparser.GetParser(bankType)
	if err != nil {
		return nil, fmt.Errorf("getting parser: %w", err)
	}

	rows, err := parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing csv: %w", err)
	}

	transactions, err := mapRowsToTransactions(rows, bankType)
	if err != nil {
		return nil, fmt.Errorf("mapping rows: %w", err)
	}

	return transactions, nil
}
