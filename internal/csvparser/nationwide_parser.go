package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

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

	for rowNum := 1; ; rowNum++ {
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

		date, err := time.Parse("02 Jan 2006", row[dateIdx])
		if err != nil {
			return nil, fmt.Errorf("row %d: parsing date '%s': %w", rowNum, row[dateIdx], err)
		}

		paidOut := row[paidOutIdx]
		paidIn := row[paidInIdx]

		amountStr := paidIn
		if paidOut != "" {
			amountStr = "-" + paidOut
		}

		amount, err := parseAmount(amountStr)
		if err != nil {
			return nil, fmt.Errorf("row %d: parsing amount '%s': %w", rowNum, amountStr, err)
		}

		transactions = append(transactions, transaction.Transaction{
			Date:        date,
			Description: row[descriptionIdx],
			Amount:      amount,
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
