package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kushturner/finances/internal/handlers"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", handlers.Ping)

	return r
}
