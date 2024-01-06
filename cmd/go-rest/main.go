package main

/*
 * This file was last modified at 2024.01.04 15:07 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * main.go
 * $Id$
 */

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"rest_api/internal/config"
	"rest_api/internal/http-server/handlers/redirect"
	"rest_api/internal/http-server/handlers/url/save"
	"rest_api/internal/http-server/middleware/logger"
	"rest_api/internal/lib/logger/handlers/slogpretty"
	"rest_api/internal/lib/logger/sl"
	"rest_api/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// Init config: cleanenv
	cfg := config.MustLoad()

	// Init logger: slog
	log := initLogger(cfg)

	// Init storage: sqlite3
	storage := initStorage(cfg, log)

	// Init router: chi, "chi render"
	router := initRouter(log, storage)

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Run server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func initLogger(cfg *config.Config) *slog.Logger {

	log := slog.New(createLoggerHandler(cfg.Env))
	log.Debug("debug messages are enabled")
	if cfg.Env == envLocal {
		log.Debug(fmt.Sprintf("%#v", cfg))
	}
	log.Info(
		"starting go-rest-api",
		slog.String("env", cfg.Env),
		slog.String("version", "0.1"),
	)
	return log
}

func initStorage(cfg *config.Config, log *slog.Logger) *sqlite.Storage {

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	return storage
}

func initRouter(log *slog.Logger, storage *sqlite.Storage) *chi.Mux {

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	return router
}

func createLoggerHandler(env string) slog.Handler {

	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		return opts.NewPrettyHandler(os.Stdout)
	case envDev:
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	default: // envProd
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
}

/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
