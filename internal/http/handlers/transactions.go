package handlers

import (
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/kushturner/finances/internal/json"
	"github.com/kushturner/finances/internal/statements"
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
	llm     statements.LLMClient
}

func NewTransactionHandler(service transactions.Service, llm statements.LLMClient) *transactionHandler {
	return &transactionHandler{
		service: service,
		llm:     llm,
	}
}

func (h *transactionHandler) AddTransactions(w http.ResponseWriter, r *http.Request) {

	fileBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		slog.Error("Failed to read statement file", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactions, err := h.llm.ReadStatement(r.Context(), base64.StdEncoding.EncodeToString(fileBytes))

	if err != nil {
		slog.Error("LLM Failed to read statement", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.service.AddTransactions(r.Context(), transactions)

	json.Encode(w, http.StatusAccepted, "Added transactions successfully")
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
