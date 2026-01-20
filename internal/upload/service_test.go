package upload

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kushturner/finances/internal/db"
	"github.com/stretchr/testify/assert"
)

type mockQuerier struct {
	createTransactionsBatchFunc func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error)
}

func (m *mockQuerier) CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) CreateTransactionsBatch(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
	if m.createTransactionsBatchFunc != nil {
		return m.createTransactionsBatchFunc(ctx, arg)
	}
	return 0, nil
}

func (m *mockQuerier) DeleteTransaction(ctx context.Context, id int32) error {
	return nil
}

func (m *mockQuerier) GetTransaction(ctx context.Context, id int32) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) ListTransactions(ctx context.Context) ([]db.Transaction, error) {
	return nil, nil
}

func (m *mockQuerier) UpdateTransaction(ctx context.Context, arg db.UpdateTransactionParams) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func TestService_UploadTransactions_Success_Nationwide(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			assert.Equal(t, 2, len(arg))
			assert.Equal(t, "nationwide", arg[0].Bank)
			return int64(len(arg)), nil
		},
	}

	service := NewService(mockQuerier)

	csv := `Account Name: Current Account
Balance: £1,000.00
Available Balance: £950.00

Date,Transaction type,Description,Paid out,Paid in,Balance
15 Jan 2026,DEB,Coffee Shop,£5.00,,£995.00
14 Jan 2026,DEB,Grocery Store,£50.00,,£945.00`

	count, err := service.UploadTransactions(context.Background(), "nationwide", strings.NewReader(csv))

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestService_UploadTransactions_Success_Amex(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			assert.Equal(t, 2, len(arg))
			assert.Equal(t, "amex", arg[0].Bank)
			return int64(len(arg)), nil
		},
	}

	service := NewService(mockQuerier)

	csv := `Date,Description,Amount,Category
15/01/2026,Restaurant,-25.50,Entertainment-Restaurants
14/01/2026,Grocery Store,-50.00,Shopping-Groceries`

	count, err := service.UploadTransactions(context.Background(), "amex", strings.NewReader(csv))

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestService_UploadTransactions_InvalidBankType(t *testing.T) {
	mockQuerier := &mockQuerier{}
	service := NewService(mockQuerier)

	csv := `Date,Description,Amount
15/01/2026,Test,-10.00`

	count, err := service.UploadTransactions(context.Background(), "invalid_bank", strings.NewReader(csv))

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.ErrorIs(t, err, ErrInvalidBankType)
}

func TestService_UploadTransactions_CSVParsingError(t *testing.T) {
	mockQuerier := &mockQuerier{}
	service := NewService(mockQuerier)

	csv := `Date,Description,Amount
invalid_date,Test,-10.00`

	count, err := service.UploadTransactions(context.Background(), "amex", strings.NewReader(csv))

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.ErrorIs(t, err, ErrParseFailure)
}

func TestService_UploadTransactions_DatabaseError(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			return 0, errors.New("database connection failed")
		},
	}

	service := NewService(mockQuerier)

	csv := `Date,Description,Amount
15/01/2026,Test,-10.00`

	count, err := service.UploadTransactions(context.Background(), "amex", strings.NewReader(csv))

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.ErrorIs(t, err, ErrDatabaseFailure)
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestService_UploadTransactions_EmptyCSV(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			assert.Equal(t, 0, len(arg))
			return 0, nil
		},
	}

	service := NewService(mockQuerier)

	csv := `Date,Description,Amount`

	count, err := service.UploadTransactions(context.Background(), "amex", strings.NewReader(csv))

	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
