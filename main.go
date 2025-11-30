package main

import (
	"finance1/internal/config"
	"finance1/internal/lib/logger/sl"
	"finance1/internal/storage"
	"os"

	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	evnDev   = "development"
	envProd  = "prod"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting server")
	log.Debug("debug logging enabled", slog.String("env", cfg.Env))

	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("error creating storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case evnDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log

}
