package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/transaction"
)

func NewListTransactionsHandler(querier db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbTransactions, err := querier.ListTransactions(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
			return
		}

		responses := make([]TransactionResponse, 0, len(dbTransactions))
		for _, dbTx := range dbTransactions {
			domainTx := transaction.TransactionFromDB(dbTx)
			responses = append(responses, FromTransaction(domainTx))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
