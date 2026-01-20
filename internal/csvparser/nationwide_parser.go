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

type NationwideParser struct{}

func (p *NationwideParser) Parse(r io.Reader) ([]transaction.Transaction, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading first row: %w", err)
	}
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading second row: %w", err)
	}
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading third row: %w", err)
	}

	var headers []string
	for {
		row, err := reader.Read()
		if err != nil {
			return nil, fmt.Errorf("reading header row: %w", err)
		}
		if len(row) > 0 && row[0] != "" {
			headers = row
			break
		}
	}

	dateIdx := findColumnIndex(headers, "Date")
	descriptionIdx := findColumnIndex(headers, "Description")
	paidOutIdx := findColumnIndex(headers, "Paid out")
	paidInIdx := findColumnIndex(headers, "Paid in")

	if dateIdx == -1 || descriptionIdx == -1 || paidOutIdx == -1 || paidInIdx == -1 {
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

		if len(row) <= max(dateIdx, descriptionIdx, paidOutIdx, paidInIdx) {
			return nil, fmt.Errorf("row has fewer columns than expected")
		}

		dateStr := row[dateIdx]
		description := row[descriptionIdx]
		paidOut := row[paidOutIdx]
		paidIn := row[paidInIdx]

		parsedDate, err := time.Parse("02 Jan 2006", dateStr)
		if err != nil {
			return nil, fmt.Errorf("parsing date '%s': %w", dateStr, err)
		}

		var amountCents int64
		if paidOut != "" {
			amount, err := parseAmount(paidOut)
			if err != nil {
				return nil, fmt.Errorf("parsing paid out amount: %w", err)
			}
			amountCents = -amount
		} else if paidIn != "" {
			amount, err := parseAmount(paidIn)
			if err != nil {
				return nil, fmt.Errorf("parsing paid in amount: %w", err)
			}
			amountCents = amount
		}

		transactions = append(transactions, transaction.Transaction{
			Date:        parsedDate,
			Description: description,
			Amount:      money.New(amountCents, "GBP"),
			Bank:        "nationwide",
			Category:    nil,
		})
	}

	return transactions, nil
}

func findColumnIndex(headers []string, columnName string) int {
	for i, header := range headers {
		if header == columnName {
			return i
		}
	}
	return -1
}

func parseAmount(amountStr string) (int64, error) {
	cleaned := strings.ReplaceAll(amountStr, "£", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	amount, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, err
	}

	return int64(amount * 100), nil
}

func max(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	maxVal := nums[0]
	for _, n := range nums[1:] {
		if n > maxVal {
			maxVal = n
		}
	}
	return maxVal
}
