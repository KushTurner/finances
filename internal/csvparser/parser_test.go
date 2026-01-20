package csvparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParser_Nationwide(t *testing.T) {
	parser, err := GetParser("nationwide")

	assert.NoError(t, err)
	assert.NotNil(t, parser)
	assert.IsType(t, &NationwideParser{}, parser)
}

func TestGetParser_Amex(t *testing.T) {
	parser, err := GetParser("amex")

	assert.NoError(t, err)
	assert.NotNil(t, parser)
	assert.IsType(t, &AmexParser{}, parser)
}

func TestGetParser_InvalidBankType(t *testing.T) {
	parser, err := GetParser("invalid")

	assert.Error(t, err)
	assert.Nil(t, parser)
	assert.Contains(t, err.Error(), "unsupported bank type")
}
