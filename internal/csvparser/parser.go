package csvparser

import (
	"fmt"
	"io"
	"strings"
)

type Parser interface {
	Parse(r io.Reader) ([]TransactionRow, error)
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

type TransactionRow struct {
	Date        string
	Description string
	Amount      string
	Bank        string
	Category    string
}
