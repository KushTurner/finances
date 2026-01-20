package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/server"
	"github.com/kushturner/finances/migrations"
)

func main() {
	connStr := os.Getenv("FINANCES_DATABASE_URL")
	if connStr == "" {
		log.Fatal("FINANCES_DATABASE_URL environment variable is not set")
	}

	if err := db.RunMigrations(connStr, migrations.FS); err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations completed")

	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	log.Println("Connected to database")

	querier := db.New(conn)
	r := server.NewRouter(querier)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
