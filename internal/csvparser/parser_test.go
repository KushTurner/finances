package csvparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Parse_InvalidBankType(t *testing.T) {
	svc := NewService()
	transactions, err := svc.Parse(strings.NewReader(""), "invalid")

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Contains(t, err.Error(), "unsupported bank type")
}
