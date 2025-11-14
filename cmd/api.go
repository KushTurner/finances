package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func startApi(h http.Handler, port string) {
	msg := fmt.Sprintf("starting server on port %s", port)
	slog.Info(msg)

	srv := &http.Server{
		Addr:    port,
		Handler: h,
	}

	if err := srv.ListenAndServe(); err != nil {
		msg := fmt.Sprintf("port %s is already being used", port)
		slog.Error(msg)
	}
}
