package transaction

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kushturner/finances/internal/db"
	"github.com/stretchr/testify/assert"
)

type mockQuerier struct {
	transactions                []db.Transaction
	err                         error
	createTransactionsBatchFunc func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error)
}

type mockParser struct {
	parseFunc func(r io.Reader) ([]Transaction, error)
}

func (m *mockParser) Parse(r io.Reader) ([]Transaction, error) {
	if m.parseFunc != nil {
		return m.parseFunc(r)
	}
	return nil, nil
}

func (m *mockQuerier) ListTransactions(ctx context.Context) ([]db.Transaction, error) {
	return m.transactions, m.err
}

func (m *mockQuerier) GetTransaction(ctx context.Context, id int32) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) UpdateTransaction(ctx context.Context, arg db.UpdateTransactionParams) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func (m *mockQuerier) DeleteTransaction(ctx context.Context, id int32) error {
	return nil
}

func (m *mockQuerier) CreateTransactionsBatch(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
	if m.createTransactionsBatchFunc != nil {
		return m.createTransactionsBatchFunc(ctx, arg)
	}
	return 0, nil
}

func (m *mockQuerier) WithTx(tx interface{}) db.Querier {
	return m
}

func TestService_GetAllTransactions_EmptyList(t *testing.T) {
	mock := &mockQuerier{
		transactions: []db.Transaction{},
		err:          nil,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, transactions)
}

func TestService_GetAllTransactions_MultipleTransactions(t *testing.T) {
	category := "groceries"
	mock := &mockQuerier{
		transactions: []db.Transaction{
			{
				ID:          1,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "Grocery store",
				Amount:      5000,
				Currency:    "USD",
				Bank:        "Chase",
				Category:    pgtype.Text{String: category, Valid: true},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
			{
				ID:          2,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "Coffee shop",
				Amount:      500,
				Currency:    "GBP",
				Bank:        "Barclays",
				Category:    pgtype.Text{String: "", Valid: false},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		},
		err: nil,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	assert.Equal(t, int32(1), transactions[0].ID)
	assert.Equal(t, "Grocery store", transactions[0].Description)
	assert.Equal(t, int64(5000), transactions[0].Amount.Amount())
	assert.Equal(t, "USD", transactions[0].Amount.Currency().Code)
	assert.Equal(t, "Chase", transactions[0].Bank)
	assert.NotNil(t, transactions[0].Category)
	assert.Equal(t, "groceries", *transactions[0].Category)

	assert.Equal(t, int32(2), transactions[1].ID)
	assert.Equal(t, "Coffee shop", transactions[1].Description)
	assert.Equal(t, int64(500), transactions[1].Amount.Amount())
	assert.Equal(t, "GBP", transactions[1].Amount.Currency().Code)
	assert.Equal(t, "Barclays", transactions[1].Bank)
	assert.Nil(t, transactions[1].Category)
}

func TestService_GetAllTransactions_DatabaseError(t *testing.T) {
	mock := &mockQuerier{
		transactions: nil,
		err:          assert.AnError,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Equal(t, assert.AnError, err)
}

func TestService_GetAllTransactions_CategoryHandling(t *testing.T) {
	category := "transport"
	mock := &mockQuerier{
		transactions: []db.Transaction{
			{
				ID:          1,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "With category",
				Amount:      1000,
				Currency:    "USD",
				Bank:        "Test Bank",
				Category:    pgtype.Text{String: category, Valid: true},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
			{
				ID:          2,
				Date:        pgtype.Date{Time: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), Valid: true},
				Description: "Without category",
				Amount:      2000,
				Currency:    "USD",
				Bank:        "Test Bank",
				Category:    pgtype.Text{String: "", Valid: false},
				CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		},
		err: nil,
	}

	service := NewService(mock)
	transactions, err := service.GetAllTransactions(context.Background())

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	assert.NotNil(t, transactions[0].Category)
	assert.Equal(t, "transport", *transactions[0].Category)

	assert.Nil(t, transactions[1].Category)
}

func TestService_ImportFromCSV_Success_Nationwide(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			assert.Equal(t, 2, len(arg))
			assert.Equal(t, "nationwide", arg[0].Bank)
			return int64(len(arg)), nil
		},
	}

	parser := &mockParser{
		parseFunc: func(r io.Reader) ([]Transaction, error) {
			return []Transaction{
				{Bank: "nationwide", Description: "Coffee Shop", Amount: money.New(500, "GBP")},
				{Bank: "nationwide", Description: "Grocery Store", Amount: money.New(5000, "GBP")},
			}, nil
		},
	}

	service := NewService(mockQuerier)
	count, err := service.ImportFromCSV(context.Background(), parser, strings.NewReader(""))

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestService_ImportFromCSV_Success_Amex(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			assert.Equal(t, 2, len(arg))
			assert.Equal(t, "amex", arg[0].Bank)
			return int64(len(arg)), nil
		},
	}

	parser := &mockParser{
		parseFunc: func(r io.Reader) ([]Transaction, error) {
			return []Transaction{
				{Bank: "amex", Description: "Restaurant", Amount: money.New(-2550, "GBP")},
				{Bank: "amex", Description: "Grocery Store", Amount: money.New(-5000, "GBP")},
			}, nil
		},
	}

	service := NewService(mockQuerier)
	count, err := service.ImportFromCSV(context.Background(), parser, strings.NewReader(""))

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestService_ImportFromCSV_CSVParsingError(t *testing.T) {
	mockQuerier := &mockQuerier{}

	parser := &mockParser{
		parseFunc: func(r io.Reader) ([]Transaction, error) {
			return nil, errors.New("invalid date format")
		},
	}

	service := NewService(mockQuerier)
	count, err := service.ImportFromCSV(context.Background(), parser, strings.NewReader(""))

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.ErrorIs(t, err, ErrParseFailure)
}

func TestService_ImportFromCSV_DatabaseError(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			return 0, errors.New("database connection failed")
		},
	}

	parser := &mockParser{
		parseFunc: func(r io.Reader) ([]Transaction, error) {
			return []Transaction{
				{Bank: "amex", Description: "Test", Amount: money.New(-1000, "GBP")},
			}, nil
		},
	}

	service := NewService(mockQuerier)
	count, err := service.ImportFromCSV(context.Background(), parser, strings.NewReader(""))

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.ErrorIs(t, err, ErrDatabaseFailure)
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestService_ImportFromCSV_EmptyCSV(t *testing.T) {
	mockQuerier := &mockQuerier{
		createTransactionsBatchFunc: func(ctx context.Context, arg []db.CreateTransactionsBatchParams) (int64, error) {
			assert.Equal(t, 0, len(arg))
			return 0, nil
		},
	}

	parser := &mockParser{
		parseFunc: func(r io.Reader) ([]Transaction, error) {
			return []Transaction{}, nil
		},
	}

	service := NewService(mockQuerier)
	count, err := service.ImportFromCSV(context.Background(), parser, strings.NewReader(""))

	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
