package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sa3akash/first-rest-api-go/internal/config"
)

func main() {
	slog.Info("welcome to rest api.")

	// load config
	cfg := config.MustLoad()

	// database

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcode to rest api"))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	slog.Info("Server start", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cencel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
