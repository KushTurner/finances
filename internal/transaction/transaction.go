package transaction

import (
	"io"
	"time"

	"github.com/Rhymond/go-money"
)

type Transaction struct {
	ID          int32
	Date        time.Time
	Description string
	Amount      *money.Money
	Bank        string
	Category    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Parser interface {
	Parse(r io.Reader) ([]Transaction, error)
}
