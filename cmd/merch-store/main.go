package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justcgh9/merch_store/internal/config"
	// "github.com/justcgh9/merch_store/internal/http-server/handlers/auth"
	mySlog "github.com/justcgh9/merch_store/internal/log"
	"github.com/justcgh9/merch_store/internal/storage/postgres"
)

func main() {

	cfg := config.MustLoad()

	log := mySlog.SetupLogger(cfg.Env)

	log.Info("starting merch store", slog.String("env", cfg.Env))

	storage := postgres.New(cfg.StoragePath)
	defer func() {
		err := storage.Close()
		if err != nil {
			log.Error("error closing db connection", slog.String("err", err.Error()))
		}
	}()

	log.Info("connected to postgres")


	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Post("/api/auth", auth.New(log, 1))

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout: cfg.IddleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("err", err.Error()))

		return
	}

	log.Info("server stopped")
}
