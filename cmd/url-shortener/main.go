package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"url-shortener/internal/config"
	"url-shortener/internal/storage/sqlite"

	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/sl"

	"github.com/MatusOllah/slogcolor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting url-shortener service")
	log.Debug("debug messages enabled", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("host", cfg.HTTPServer.Host))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Host + ":" + strconv.Itoa(cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("service stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	var opts slogcolor.Options = *slogcolor.DefaultOptions

	switch env {
	case envLocal:
		opts.Level = slog.LevelDebug
		log = slog.New(slogcolor.NewHandler(os.Stderr, &opts))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		opts.Level = slog.LevelInfo
		log = slog.New(slogcolor.NewHandler(os.Stderr, &opts))
	}

	return log
}
