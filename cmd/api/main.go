package main

import (
	"log"
	"net/http"

	"github.com/kushturner/finances/internal/server"
)

func main() {
	r := server.NewRouter()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
