package csvparser

import (
	"os"
	"strings"
	"testing"

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

	assert.Equal(t, "15 Jan 2026", transactions[0].Date)
	assert.Equal(t, "TEST PAYEE", transactions[0].Description)
	assert.Equal(t, "-£50.00", transactions[0].Amount)
	assert.Equal(t, "nationwide", transactions[0].Bank)
	assert.Equal(t, "", transactions[0].Category)

	assert.Equal(t, "TEST MERCHANT LONDON GB APPLEPAY 1234", transactions[1].Description)
	assert.Equal(t, "-£10.50", transactions[1].Amount)

	assert.Equal(t, "TEST STANDING ORDER", transactions[2].Description)
	assert.Equal(t, "-£100.00", transactions[2].Amount)

	assert.Equal(t, "TEST UTILITY COMPANY", transactions[3].Description)
	assert.Equal(t, "-£75.25", transactions[3].Amount)

	assert.Equal(t, "SALARY PAYMENT", transactions[4].Description)
	assert.Equal(t, "£2000.00", transactions[4].Amount)
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
