package csvparser

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAmexParser_Parse_ValidCSV(t *testing.T) {
	file, err := os.Open("testdata/amex_sample.csv")
	assert.NoError(t, err)
	defer file.Close()

	parser := &AmexParser{}
	transactions, err := parser.Parse(file)

	assert.NoError(t, err)
	assert.Equal(t, 4, len(transactions))

	assert.Equal(t, "TEST RESTAURANT LONDON", transactions[0].Description)
	assert.Equal(t, int64(2550), transactions[0].Amount.Amount())
	assert.Equal(t, "GBP", transactions[0].Amount.Currency().Code)
	assert.Equal(t, "amex", transactions[0].Bank)
	assert.NotNil(t, transactions[0].Category)
	assert.Equal(t, "Entertainment-Restaurants", *transactions[0].Category)
	expectedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedDate, transactions[0].Date)

	assert.Equal(t, "PAYMENT RECEIVED - THANK YOU", transactions[1].Description)
	assert.Equal(t, int64(-10000), transactions[1].Amount.Amount())
	assert.Nil(t, transactions[1].Category)

	assert.Equal(t, "TEST COFFEE SHOP LONDON", transactions[2].Description)
	assert.Equal(t, int64(475), transactions[2].Amount.Amount())
	assert.NotNil(t, transactions[2].Category)
	assert.Equal(t, "Entertainment-Bars & Cafés", *transactions[2].Category)

	assert.Equal(t, "TEST SUPERMARKET LONDON", transactions[3].Description)
	assert.Equal(t, int64(4520), transactions[3].Amount.Amount())
	assert.NotNil(t, transactions[3].Category)
	assert.Equal(t, "Shopping-Groceries", *transactions[3].Category)
}

func TestAmexParser_Parse_InvalidDateFormat(t *testing.T) {
	csvData := `Date,Description,Card Member,Account #,Amount,Extended Details,Appears On Your Statement As,Address,Town/City,Postcode,Country,Reference,Category
invalid date,TEST,MR TEST,-12345,25.50,,,,,,,AT123,Test Category
`
	parser := &AmexParser{}
	_, err := parser.Parse(strings.NewReader(csvData))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing date")
}

func TestAmexParser_Parse_InvalidAmount(t *testing.T) {
	csvData := `Date,Description,Card Member,Account #,Amount,Extended Details,Appears On Your Statement As,Address,Town/City,Postcode,Country,Reference,Category
15/01/2026,TEST,MR TEST,-12345,invalid,,,,,,,AT123,Test Category
`
	parser := &AmexParser{}
	_, err := parser.Parse(strings.NewReader(csvData))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing amount")
}
