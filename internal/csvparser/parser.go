package csvparser

import (
	"fmt"
	"io"
	"strings"

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
	parser, err := getParser(bankType)
	if err != nil {
		return nil, err
	}
	return parser.Parse(r)
}

type parser interface {
	Parse(r io.Reader) ([]transaction.Transaction, error)
}

func getParser(bankType string) (parser, error) {
	switch strings.ToLower(bankType) {
	case "nationwide":
		return &NationwideParser{}, nil
	case "amex":
		return &AmexParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported bank type: %s", bankType)
	}
}
