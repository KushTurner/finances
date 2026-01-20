package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/kushturner/finances/internal/transaction"
)

type AmexParser struct{}

func (p *AmexParser) Parse(r io.Reader) ([]transaction.Transaction, error) {
	reader := csv.NewReader(r)

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header row: %w", err)
	}

	dateIdx := findColumnIndex(headers, "Date")
	descriptionIdx := findColumnIndex(headers, "Description")
	amountIdx := findColumnIndex(headers, "Amount")
	categoryIdx := findColumnIndex(headers, "Category")

	if dateIdx == -1 || descriptionIdx == -1 || amountIdx == -1 {
		return nil, fmt.Errorf("required column not found in CSV headers")
	}

	var transactions []transaction.Transaction

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading data row: %w", err)
		}

		if len(row) <= max(dateIdx, descriptionIdx, amountIdx) {
			return nil, fmt.Errorf("row has fewer columns than expected")
		}

		dateStr := row[dateIdx]
		description := row[descriptionIdx]
		amountStr := row[amountIdx]

		parsedDate, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing date '%s': %w", dateStr, err)
		}

		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing amount '%s': %w", amountStr, err)
		}
		amountCents := int64(amountFloat * 100)

		var category *string
		if categoryIdx != -1 && len(row) > categoryIdx && row[categoryIdx] != "" {
			trimmed := strings.TrimSpace(row[categoryIdx])
			if trimmed != "" {
				category = &trimmed
			}
		}

		transactions = append(transactions, transaction.Transaction{
			Date:        parsedDate,
			Description: description,
			Amount:      money.New(amountCents, "GBP"),
			Bank:        "amex",
			Category:    category,
		})
	}

	return transactions, nil
}
