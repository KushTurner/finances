package transaction

import "errors"

var (
	ErrInvalidBankType = errors.New("invalid bank type")
	ErrParseFailure    = errors.New("parse failure")
	ErrDatabaseFailure = errors.New("database failure")
)
