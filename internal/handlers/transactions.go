package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kushturner/finances/internal/transaction"
)

func NewListTransactionsHandler(transactionService transaction.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := transactionService.GetAllTransactions(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
			return
		}

		responses := make([]TransactionResponse, 0, len(transactions))
		for _, tx := range transactions {
			responses = append(responses, FromTransaction(tx))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
