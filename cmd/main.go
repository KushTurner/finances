package main

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kushturner/finances/internal/database/psql"
	repo "github.com/kushturner/finances/internal/database/psql/sqlc"
	"github.com/kushturner/finances/internal/env"
	"github.com/kushturner/finances/internal/http/handlers"
	"github.com/kushturner/finances/internal/statements"
	"github.com/kushturner/finances/internal/transactions"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	logger := slog.Default()
	slog.SetDefault(logger)

	conn, _ := pgxpool.New(ctx, env.Get("", "postgres://postgres:postgres@localhost:5401/postgres"))
	defer conn.Close()

	llm := statements.NewLLMClient(env.Get("OPENAPI_KEY", ""))

	goose.SetBaseFS(psql.Migrations)

	if err := goose.SetDialect(env.Get("", "postgres")); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(conn)
	defer db.Close()

	if err := goose.Up(db, env.Get("", "migrations")); err != nil {
		panic(err)
	}

	txSvc := transactions.NewService(repo.New(conn))
	txHandler := handlers.NewTransactionHandler(txSvc, *llm)
	r.Get("/transactions", txHandler.ListTransactions)
	r.Post("/statement", txHandler.AddTransactions)

	startApi(r, ":3000")
}
