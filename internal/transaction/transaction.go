package transaction

import "time"

type Transaction struct {
	ID          int32
	Date        time.Time
	Description string
	Amount      int64
	Currency    string
	Bank        string
	Category    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
