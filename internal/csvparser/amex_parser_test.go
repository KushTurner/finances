package csvparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAmount_UTF8PoundSign(t *testing.T) {
	result, err := parseAmount("£100.00")

	assert.NoError(t, err)
	assert.Equal(t, int64(10000), result.Amount())
}

func TestParseAmount_ISO88591PoundSign(t *testing.T) {
	result, err := parseAmount("\xa3100.00")

	assert.NoError(t, err)
	assert.Equal(t, int64(10000), result.Amount())
}

func TestParseAmount_NegativeWithUTF8Pound(t *testing.T) {
	result, err := parseAmount("-£100.00")

	assert.NoError(t, err)
	assert.Equal(t, int64(-10000), result.Amount())
}

func TestParseAmount_NegativeWithISO88591Pound(t *testing.T) {
	result, err := parseAmount("-\xa3100.00")

	assert.NoError(t, err)
	assert.Equal(t, int64(-10000), result.Amount())
}

func TestParseAmount_WithCommas(t *testing.T) {
	result, err := parseAmount("£1,234.56")

	assert.NoError(t, err)
	assert.Equal(t, int64(123456), result.Amount())
}

func TestParseAmount_WithoutCurrencySymbol(t *testing.T) {
	result, err := parseAmount("50.75")

	assert.NoError(t, err)
	assert.Equal(t, int64(5075), result.Amount())
}

func TestParseAmount_EmptyString(t *testing.T) {
	_, err := parseAmount("")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty amount string")
}

func TestParseAmount_OnlyCurrencySymbol(t *testing.T) {
	_, err := parseAmount("£")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty amount string")
}

func TestParseAmount_InvalidFormat(t *testing.T) {
	_, err := parseAmount("abc")

	assert.Error(t, err)
}

func TestParseAmount_MultipleDecimalPoints(t *testing.T) {
	_, err := parseAmount("£10.50.25")

	assert.Error(t, err)
}

func TestParseAmount_Whitespace(t *testing.T) {
	result, err := parseAmount("  £100.00  ")

	assert.NoError(t, err)
	assert.Equal(t, int64(10000), result.Amount())
}

func TestParseAmount_SmallAmount(t *testing.T) {
	result, err := parseAmount("£0.01")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Amount())
}

func TestParseAmount_Zero(t *testing.T) {
	result, err := parseAmount("£0.00")

	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.Amount())
}
