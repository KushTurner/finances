package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
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

	for rowNum := 1; ; rowNum++ {
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

		date, err := time.Parse("02/01/2006", row[dateIdx])
		if err != nil {
			return nil, fmt.Errorf("row %d: parsing date '%s': %w", rowNum, row[dateIdx], err)
		}

		amount, err := parseAmount(row[amountIdx])
		if err != nil {
			return nil, fmt.Errorf("row %d: parsing amount '%s': %w", rowNum, row[amountIdx], err)
		}

		var category *string
		if categoryIdx != -1 && len(row) > categoryIdx {
			cat := strings.TrimSpace(row[categoryIdx])
			if cat != "" {
				category = &cat
			}
		}

		transactions = append(transactions, transaction.Transaction{
			Date:        date,
			Description: row[descriptionIdx],
			Amount:      amount,
			Bank:        "American Express",
			Category:    category,
		})
	}

	return transactions, nil
}

func parseAmount(amountStr string) (*money.Money, error) {
	re := regexp.MustCompile(`[^0-9.-]`)
	cleaned := re.ReplaceAllString(amountStr, "")
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return nil, fmt.Errorf("empty amount string")
	}

	amountFloat, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return nil, err
	}

	amountCents := int64(amountFloat * 100)
	return money.New(amountCents, "GBP"), nil
}
