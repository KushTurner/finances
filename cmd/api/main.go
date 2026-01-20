package main

import (
	"log"
	"net/http"

	"github.com/kushturner/finances/internal/db"
	"github.com/kushturner/finances/internal/server"
	"github.com/kushturner/finances/migrations"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connected to database")

	if err := db.RunMigrations(conn, migrations.FS); err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations completed")

	r := server.NewRouter()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
