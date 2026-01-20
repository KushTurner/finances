package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/handlers"
	"github.com/kushturner/finances/internal/upload"
)

func NewRouter(querier db.Querier) *chi.Mux {
	r := chi.NewRouter()

	uploadService := upload.NewService(querier)

	r.Get("/ping", handlers.Ping)
	r.Get("/transactions", handlers.NewListTransactionsHandler(querier))
	r.Post("/transactions/upload", handlers.NewUploadTransactionsHandler(uploadService))

	return r
}
