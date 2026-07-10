package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucjosin/qorv.in/internal/api"
	"github.com/lucjosin/qorv.in/internal/slogx"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerConfig struct {
	Port           string `env:"PORT" envDefault:":4512"`
	URL            string `env:"URL"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`
}

type Config struct {
	Server ServerConfig `envPrefix:"SERVER_"`
}

func main() {
	// setup signal handling for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx, log, err := slogx.NewWithContext(ctx)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Panic(err)
	}
	if cfg.Server.URL == "" {
		cfg.Server.URL = "http://0.0.0.0" + cfg.Server.Port
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(middleware.AllowContentType("application/json"))

	r.Route("/api", func(r chi.Router) {
		// public routes
		api.NewHandler().RegisterRoutes(r)
	})

	errLog := slog.NewLogLogger(log.Handler(), slog.LevelError)
	server := http.Server{
		Handler:      r,
		Addr:         cfg.Server.Port,
		ErrorLog:     errLog,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Info("server listening and serving on " + cfg.Server.Port)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()

	// interrupt signal
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Panic(err)
	}

	log.Info("server stopped")
}
