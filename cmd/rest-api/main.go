package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BinayRajbanshi/go-rest-API/database/sqlite"
	config "github.com/BinayRajbanshi/go-rest-API/internal"
	"github.com/BinayRajbanshi/go-rest-API/internal/controllers/user"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	db, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database initialized", slog.String("Env: ", cfg.Env))
	// setup router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/users", user.New(db))
	mux.HandleFunc("GET /api/v1/users", user.GetAll(db))
	mux.HandleFunc("DELETE /api/v1/users/{id}", user.Delete(db))

	// setup server
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	slog.Info("Server started: ", slog.String("Address", cfg.Address))

	done := make(chan os.Signal, 1) // The buffer ensures that if multiple signals are received rapidly, none are lost.

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server is closed")
		} else if err != nil {
			slog.Error("Error listening to the server: ")
		}
	}()

	<-done

	slog.Info("Shutting down the server")
	ctx := context.Background()
	ctx, cancelContext := context.WithTimeout(ctx, 5*time.Second) //ctx.Done() channel is called automatically after 5 seconds. I can also check if context.DeadlineExceeded or context.Canceled
	defer cancelContext()

	err = server.Shutdown(ctx) //if the provided context expires before the shutdown is complete, Shutdown returns the context's error, otherwise it returns any error returned from closing the [Server]'s underlying Listener(s).
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successful")
}
