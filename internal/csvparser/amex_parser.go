package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type AmexParser struct{}

func (p *AmexParser) Parse(r io.Reader) ([]TransactionRow, error) {
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

	var transactions []TransactionRow

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

		category := ""
		if categoryIdx != -1 && len(row) > categoryIdx {
			category = strings.TrimSpace(row[categoryIdx])
		}

		transactions = append(transactions, TransactionRow{
			Date:        row[dateIdx],
			Description: row[descriptionIdx],
			Amount:      row[amountIdx],
			Bank:        "amex",
			Category:    category,
		})
	}

	return transactions, nil
}
