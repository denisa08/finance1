package main

import (
	"errors"
	"finance1/internal/config"
	"finance1/internal/http-server/handlers/redirect"
	"finance1/internal/http-server/handlers/url/save"
	mwLogger "finance1/internal/http-server/middleware/logger"
	"finance1/internal/http-server/middleware/logger/handlers/slogpretty"
	"time"

	"finance1/internal/lib/logger/sl"
	"finance1/internal/storage"
	"net/http"
	"os"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	evnDev   = "development"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env))

	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("error creating storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat)
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))

	})

	router.Get("/{alias}", redirect.New(log, storage))
	//router.Delete("/{alias}", redirect.New(log, storage))

	addr := cfg.HTTPServer.Address
	log.Info("starting http server", slog.String("addr", addr))

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("http server stopped", sl.Err(err))
		os.Exit(1)
	}

}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return setupPrettySlog()
	case evnDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		// чтоб не получить nil-логгер и потом весело паниковать
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
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
