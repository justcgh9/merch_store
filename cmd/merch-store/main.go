package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justcgh9/merch_store/internal/config"
	"github.com/justcgh9/merch_store/internal/services/user"

	"github.com/justcgh9/merch_store/internal/http-server/handlers/auth"
	authMiddleware "github.com/justcgh9/merch_store/internal/http-server/middleware/auth"
	mySlog "github.com/justcgh9/merch_store/internal/log"
	"github.com/justcgh9/merch_store/internal/storage/postgres"
)

func main() {

	cfg := config.MustLoad()

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		ps := flag.String("jwt-secret", "", "jwt secret")
		flag.Parse()

		jwtSecret = *ps
	}

	if jwtSecret == "" {
		log.Fatalf("no jwt secret specified")
	}

	log := mySlog.SetupLogger(cfg.Env)

	log.Info("starting merch store", slog.String("env", cfg.Env))

	log.Info("starting postgresql connection")

	storage := postgres.New(cfg.StoragePath, cfg.Timeout)
	defer func() {
		err := storage.Close()
		if err != nil {
			log.Error("error closing db connection", slog.String("err", err.Error()))
		}
	}()

	log.Info("connected to postgres")

	userService := user.New(log, jwtSecret, storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	middleware := authMiddleware.New(log, userService)
	_ = middleware

	router.Post("/api/auth", auth.New(log, userService))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IddleTimeout,
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
