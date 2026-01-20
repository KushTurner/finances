package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kushturner/finances/internal/upload"
	"github.com/stretchr/testify/assert"
)

type mockUploadService struct {
	uploadFunc func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error)
}

func (m *mockUploadService) UploadTransactions(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, bankType, csvFile)
	}
	return 0, nil
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
	csvContent := `"Account Name:","Debit ****12345"
"Account Balance:","£1234.56"
"Available Balance: ","£1234.56"

"Date","Transaction type","Description","Paid out","Paid in","Balance"
"15 Jan 2026","Payment to","TEST PAYEE","£50.00","","£1184.56"
"14 Jan 2026","Direct Credit","SALARY","","£2000.00","£3184.56"
`

	mock := &mockUploadService{
		uploadFunc: func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
			assert.Equal(t, "nationwide", bankType)
			return 2, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "nationwide")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UploadResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully uploaded 2 transactions", response.Message)
}

func TestUploadTransactionsHandler_Success_Amex(t *testing.T) {
	csvContent := `Date,Description,Card Member,Account #,Amount,Extended Details,Appears On Your Statement As,Address,Town/City,Postcode,Country,Reference,Category
15/01/2026,TEST RESTAURANT,MR TEST,-12345,25.50,GOODS,TEST RESTAURANT,123 ST,LONDON,SW1A 1AA,UK,'AT123',Entertainment
14/01/2026,PAYMENT RECEIVED,MR TEST,-12345,-100.00,,PAYMENT RECEIVED,,,,,'AT456',
`

	mock := &mockUploadService{
		uploadFunc: func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
			assert.Equal(t, "amex", bankType)
			return 2, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "amex")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UploadResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully uploaded 2 transactions", response.Message)
}

func TestUploadTransactionsHandler_InvalidBankType(t *testing.T) {
	csvContent := "some,csv,data\n"

	mock := &mockUploadService{
		uploadFunc: func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
			return 0, fmt.Errorf("%w: unsupported bank type", upload.ErrInvalidBankType)
		},
	}

	req := createMultipartRequest(t, csvContent, "invalid")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Upload failed")
}

func TestUploadTransactionsHandler_MissingFile(t *testing.T) {
	mock := &mockUploadService{}

	req := httptest.NewRequest(http.MethodPost, "/transactions/upload?bank=nationwide", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadTransactionsHandler_CSVParsingError(t *testing.T) {
	csvContent := `"Account Name:","Debit ****12345"
"Account Balance:","£1234.56"
"Available Balance: ","£1234.56"

"Date","Transaction type","Description","Paid out","Paid in","Balance"
"invalid date","Payment to","TEST","£50.00","","£1000.00"
`

	mock := &mockUploadService{
		uploadFunc: func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
			return 0, fmt.Errorf("%w: invalid date format", upload.ErrParseFailure)
		},
	}

	req := createMultipartRequest(t, csvContent, "nationwide")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Upload failed", response.Error)
}

func TestUploadTransactionsHandler_DatabaseError(t *testing.T) {
	csvContent := `"Account Name:","Debit ****12345"
"Account Balance:","£1234.56"
"Available Balance: ","£1234.56"

"Date","Transaction type","Description","Paid out","Paid in","Balance"
"15 Jan 2026","Payment to","TEST","£50.00","","£1000.00"
`

	mock := &mockUploadService{
		uploadFunc: func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
			return 0, errors.New("failed to save transactions: database connection failed")
		},
	}

	req := createMultipartRequest(t, csvContent, "nationwide")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Upload failed", response.Error)
}

func TestUploadTransactionsHandler_MissingBankParameter(t *testing.T) {
	csvContent := `"Account Name:","Debit ****12345"
"Account Balance:","£1234.56"
"Available Balance: ","£1234.56"

"Date","Transaction type","Description","Paid out","Paid in","Balance"
"15 Jan 2026","Payment to","TEST","£50.00","","£1000.00"
`

	mock := &mockUploadService{}

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

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Missing bank type parameter", response.Error)
}

func TestUploadTransactionsHandler_EmptyCSVFile(t *testing.T) {
	csvContent := ""

	mock := &mockUploadService{
		uploadFunc: func(ctx context.Context, bankType string, csvFile io.Reader) (int64, error) {
			return 0, nil
		},
	}

	req := createMultipartRequest(t, csvContent, "nationwide")
	rec := httptest.NewRecorder()

	handler := NewUploadTransactionsHandler(mock)
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response UploadResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully uploaded 0 transactions", response.Message)
}
