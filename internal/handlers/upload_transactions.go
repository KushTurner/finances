package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kushturner/finances/internal/parser"
	"github.com/kushturner/finances/internal/transaction"
)

type UploadResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func NewUploadTransactionsHandler(transactionService transaction.Service, parserService parser.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			respondWithError(w, http.StatusBadRequest, "Failed to parse multipart form", err.Error())
			return
		}

		bankType := r.URL.Query().Get("bank")
		if bankType == "" {
			respondWithError(w, http.StatusBadRequest, "Missing bank type parameter", "")
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Failed to get file from form", err.Error())
			return
		}
		defer file.Close()

		transactions, err := parserService.Parse(file, bankType)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Failed to parse CSV file", err.Error())
			return
		}

		count, err := transactionService.AddTransactions(r.Context(), transactions)
		if err != nil {
			statusCode := determineStatusCode(err)
			respondWithError(w, statusCode, "Upload failed", err.Error())
			return
		}

		respondWithSuccess(w, count)
	}
}

func determineStatusCode(err error) int {
	if errors.Is(err, transaction.ErrParseFailure) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, transaction.ErrDatabaseFailure) {
		return http.StatusInternalServerError
	}
	return http.StatusInternalServerError
}

func respondWithSuccess(w http.ResponseWriter, count int64) {
	response := UploadResponse{
		Message: fmt.Sprintf("Successfully uploaded %d transactions", count),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, errorMsg string, details string) {
	response := ErrorResponse{
		Error:   errorMsg,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
