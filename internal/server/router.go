package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kushturner/finances/internal/csvparser"
	"github.com/kushturner/finances/internal/handlers"
	"github.com/kushturner/finances/internal/transaction"
)

func NewRouter(transactionService transaction.Service, parserService csvparser.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/transactions", handlers.NewListTransactionsHandler(transactionService))
	r.Post("/transactions/upload", handlers.NewUploadTransactionsHandler(transactionService, parserService))

	return r
}
