package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/kushturner/finances/internal/csvparser"
	"github.com/kushturner/finances/internal/transaction"
)

var dateFormats = map[string]string{
	"amex":       "02/01/2006",
	"nationwide": "02 Jan 2006",
}

func mapRowsToTransactions(rows []csvparser.TransactionRow, bankType string) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, 0, len(rows))

	dateFormat, ok := dateFormats[bankType]
	if !ok {
		return nil, fmt.Errorf("unknown bank type: %s", bankType)
	}

	for i, row := range rows {
		tx, err := mapRowToTransaction(row, dateFormat)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+1, err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func mapRowToTransaction(row csvparser.TransactionRow, dateFormat string) (transaction.Transaction, error) {
	date, err := time.Parse(dateFormat, row.Date)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("parsing date '%s': %w", row.Date, err)
	}

	amount, err := parseAmount(row.Amount)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("parsing amount '%s': %w", row.Amount, err)
	}

	var category *string
	if row.Category != "" {
		category = &row.Category
	}

	return transaction.Transaction{
		Date:        date,
		Description: row.Description,
		Amount:      amount,
		Bank:        row.Bank,
		Category:    category,
	}, nil
}

func parseAmount(amountStr string) (*money.Money, error) {
	cleaned := strings.ReplaceAll(amountStr, "£", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	amountFloat, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return nil, err
	}

	amountCents := int64(amountFloat * 100)
	return money.New(amountCents, "GBP"), nil
}
