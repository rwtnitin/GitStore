package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var PORT string = "8080"

func main() {

	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("GET /health", HealthHandler)

	var server *http.Server = &http.Server{
		Addr:    net.JoinHostPort("", PORT),
		Handler: mux,
	}

	var stop chan os.Signal = make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Println("Starting server on PORT", PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server on PORT %s: %s\n", PORT, err)
		}
	}()

	var receivedSignal os.Signal = <-stop
	log.Printf("\nReceived signal: '%s'. Shutting down gracefully...\n", receivedSignal)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v\n", err)
	}
	log.Println("Server shut down successfully")
}
