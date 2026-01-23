package csvparser

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNationwideParser_Parse_ValidCSV(t *testing.T) {
	file, err := os.Open("testdata/nationwide_sample.csv")
	assert.NoError(t, err)
	defer file.Close()

	parser := &NationwideParser{}
	transactions, err := parser.Parse(file)

	assert.NoError(t, err)
	assert.Equal(t, 5, len(transactions))

	assert.Equal(t, time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), transactions[0].Date)
	assert.Equal(t, "TEST PAYEE", transactions[0].Description)
	assert.Equal(t, int64(-5000), transactions[0].Amount.Amount())
	assert.Equal(t, "Nationwide", transactions[0].Bank)
	assert.Nil(t, transactions[0].Category)

	assert.Equal(t, "TEST MERCHANT LONDON GB APPLEPAY 1234", transactions[1].Description)
	assert.Equal(t, int64(-1050), transactions[1].Amount.Amount())

	assert.Equal(t, "TEST STANDING ORDER", transactions[2].Description)
	assert.Equal(t, int64(-10000), transactions[2].Amount.Amount())

	assert.Equal(t, "TEST UTILITY COMPANY", transactions[3].Description)
	assert.Equal(t, int64(-7525), transactions[3].Amount.Amount())

	assert.Equal(t, "SALARY PAYMENT", transactions[4].Description)
	assert.Equal(t, int64(200000), transactions[4].Amount.Amount())
}

func TestNationwideParser_Parse_MissingColumns(t *testing.T) {
	csvData := `"Account Name:","Debit ****12345"
"Account Balance:","£1234.56"
"Available Balance: ","£1234.56"

"Date","Description"
"15 Jan 2026","TEST"
`
	parser := &NationwideParser{}
	_, err := parser.Parse(strings.NewReader(csvData))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required column")
}
