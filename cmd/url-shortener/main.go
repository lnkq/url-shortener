package main

import (
	"log/slog"
	"os"

	"url-shortener/internal/config"
	"url-shortener/internal/storage/sqlite"

	sl "url-shortener/internal/lib/logger/slog"

	"github.com/MatusOllah/slogcolor"
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

	_ = storage // TODO: use storage

	// TODO: init router: chi, "chi render"
	// TODO: run server
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
