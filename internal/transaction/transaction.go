package transaction

import (
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
