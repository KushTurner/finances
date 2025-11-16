package handlers

import (
	"net/http"
	"time"

	"github.com/kushturner/finances/internal/json"
	"github.com/kushturner/finances/internal/transactions"
)

type transactionResponse struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Out         string    `json:"out"`
	In          string    `json:"in"`
}

type transactionHandler struct {
	service transactions.Service
}

func NewTransactionHandler(service transactions.Service) *transactionHandler {
	return &transactionHandler{
		service: service,
	}
}

func (h *transactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	txs, err := h.service.ListTransactions(r.Context())

	if err != nil {
		return
	}

	txsResponse := make([]transactionResponse, len(txs))
	for i, tx := range txs {
		txsResponse[i] = toTransactionResponse(tx)
	}

	json.Encode(w, http.StatusOK, txsResponse)
}

func toTransactionResponse(tx transactions.Transaction) transactionResponse {
	return transactionResponse{
		Date:        tx.Date,
		Description: tx.Description,
		Out:         tx.Out.Display(),
		In:          tx.In.Display(),
	}
}
