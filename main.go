package main

import (
	"finance1/internal/config"
	mwLogger "finance1/internal/http-server/middleware/logger"
	"finance1/internal/http-server/middleware/logger/handlers/slogpretty"

	"finance1/internal/lib/logger/sl"
	"finance1/internal/storage"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	//middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	_ = storage

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case evnDev:
		//	log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
