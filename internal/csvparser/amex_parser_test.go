package csvparser

import (
	"os"
	"strings"
	"testing"

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

	assert.Equal(t, "15/01/2026", transactions[0].Date)
	assert.Equal(t, "TEST RESTAURANT LONDON", transactions[0].Description)
	assert.Equal(t, "25.50", transactions[0].Amount)
	assert.Equal(t, "amex", transactions[0].Bank)
	assert.Equal(t, "Entertainment-Restaurants", transactions[0].Category)

	assert.Equal(t, "PAYMENT RECEIVED - THANK YOU", transactions[1].Description)
	assert.Equal(t, "-100.00", transactions[1].Amount)
	assert.Equal(t, "", transactions[1].Category)

	assert.Equal(t, "TEST COFFEE SHOP LONDON", transactions[2].Description)
	assert.Equal(t, "4.75", transactions[2].Amount)
	assert.Equal(t, "Entertainment-Bars & Cafés", transactions[2].Category)

	assert.Equal(t, "TEST SUPERMARKET LONDON", transactions[3].Description)
	assert.Equal(t, "45.20", transactions[3].Amount)
	assert.Equal(t, "Shopping-Groceries", transactions[3].Category)
}

func TestAmexParser_Parse_MissingColumns(t *testing.T) {
	csvData := `Date,Description
15/01/2026,TEST
`
	parser := &AmexParser{}
	_, err := parser.Parse(strings.NewReader(csvData))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required column")
}
