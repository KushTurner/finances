package csvparser

import (
	"fmt"
	"io"
	"strings"

	"github.com/kushturner/finances/internal/transaction"
)

type Parser interface {
	Parse(r io.Reader) ([]transaction.Transaction, error)
}

func GetParser(bankType string) (Parser, error) {
	switch strings.ToLower(bankType) {
	case "nationwide":
		return &NationwideParser{}, nil
	case "amex":
		return &AmexParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported bank type: %s", bankType)
	}
}
