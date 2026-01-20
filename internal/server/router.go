package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/handlers"
)

func NewRouter(querier db.Querier) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", handlers.Ping)
	r.Get("/transactions", handlers.NewListTransactionsHandler(querier))

	return r
}
