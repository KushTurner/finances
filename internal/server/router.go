package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kushturner/finances/internal/csvparser"
	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/handlers"
	"github.com/kushturner/finances/internal/transaction"
)

func NewRouter(querier db.Querier) *chi.Mux {
	r := chi.NewRouter()

	parserService := csvparser.NewService()
	transactionService := transaction.NewService(querier)

	r.Get("/ping", handlers.Ping)
	r.Get("/transactions", handlers.NewListTransactionsHandler(transactionService))
	r.Post("/transactions/upload", handlers.NewUploadTransactionsHandler(transactionService, parserService))

	return r
}
