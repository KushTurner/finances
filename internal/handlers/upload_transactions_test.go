package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/kushturner/finances/internal/transaction"
	"github.com/stretchr/testify/assert"
)

type mockParserService struct {
	parseFunc func(r io.Reader, bankType string) ([]transaction.Transaction, error)
}

func (m *mockParserService) Parse(r io.Reader, bankType string) ([]transaction.Transaction, error) {
	if m.parseFunc != nil {
		return m.parseFunc(r, bankType)
	}
	return nil, nil
}

func createMultipartRequest(t *testing.T, csvContent string, bankType string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "statement.csv")
	assert.NoError(t, err)

	_, err = io.WriteString(part, csvContent)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/transactions/upload?bank="+bankType, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestUploadTransactionsHandler_Success_Nationwide(t *testing.T) {
	csvContent := `some csv content`

	mockTxService := &mockTransactionService{
		addTransactionsFunc: func(ctx context.Context, transactions []transaction.Transaction) (int64, error) {
			return int64(len(transactions)), nil
		},
	}

	mockParser := &mockParserService{
		parseFunc: func(r io.Reader, bankType string) ([]transaction.Transaction, error) {
			assert.Equal(t, "nationwide", bankType)
			return []transaction.Transaction{
				{Bank: "nationwide", Description: "TEST PAYEE", Amount: money.New(-5000, "GBP")},
				{Bank: "nationwide", Description: "SALARY", Amount: money.New(200000, "GBP")},
			}, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "nationwide")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UploadResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully uploaded 2 transactions", response.Message)
}

func TestUploadTransactionsHandler_Success_Amex(t *testing.T) {
	csvContent := `some csv content`

	mockTxService := &mockTransactionService{
		addTransactionsFunc: func(ctx context.Context, transactions []transaction.Transaction) (int64, error) {
			return int64(len(transactions)), nil
		},
	}

	mockParser := &mockParserService{
		parseFunc: func(r io.Reader, bankType string) ([]transaction.Transaction, error) {
			assert.Equal(t, "amex", bankType)
			return []transaction.Transaction{
				{Bank: "amex", Description: "TEST RESTAURANT", Amount: money.New(2550, "GBP")},
				{Bank: "amex", Description: "PAYMENT RECEIVED", Amount: money.New(-10000, "GBP")},
			}, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "amex")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UploadResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully uploaded 2 transactions", response.Message)
}

func TestUploadTransactionsHandler_ParsingError(t *testing.T) {
	csvContent := `some csv content`

	mockTxService := &mockTransactionService{}

	mockParser := &mockParserService{
		parseFunc: func(r io.Reader, bankType string) ([]transaction.Transaction, error) {
			return nil, errors.New("invalid date format")
		},
	}

	req := createMultipartRequest(t, csvContent, "amex")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to parse CSV file", response.Error)
}

func TestUploadTransactionsHandler_MissingFile(t *testing.T) {
	mockTxService := &mockTransactionService{}
	mockParser := &mockParserService{}

	req := httptest.NewRequest(http.MethodPost, "/transactions/upload?bank=nationwide", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadTransactionsHandler_DatabaseError(t *testing.T) {
	csvContent := `some csv content`

	mockTxService := &mockTransactionService{
		addTransactionsFunc: func(ctx context.Context, transactions []transaction.Transaction) (int64, error) {
			return 0, transaction.ErrDatabaseFailure
		},
	}

	mockParser := &mockParserService{
		parseFunc: func(r io.Reader, bankType string) ([]transaction.Transaction, error) {
			return []transaction.Transaction{
				{Bank: "amex", Description: "TEST", Amount: money.New(2550, "GBP")},
			}, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "amex")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Upload failed", response.Error)
}

func TestUploadTransactionsHandler_MissingBankParameter(t *testing.T) {
	csvContent := `some csv content`

	mockTxService := &mockTransactionService{}
	mockParser := &mockParserService{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "statement.csv")
	assert.NoError(t, err)

	_, err = io.WriteString(part, csvContent)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/transactions/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Missing bank type parameter", response.Error)
}

func TestUploadTransactionsHandler_EmptyCSVFile(t *testing.T) {
	csvContent := ""

	mockTxService := &mockTransactionService{
		addTransactionsFunc: func(ctx context.Context, transactions []transaction.Transaction) (int64, error) {
			return 0, nil
		},
	}

	mockParser := &mockParserService{
		parseFunc: func(r io.Reader, bankType string) ([]transaction.Transaction, error) {
			return []transaction.Transaction{}, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "nationwide")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mockTxService, mockParser)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UploadResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully uploaded 0 transactions", response.Message)
}
