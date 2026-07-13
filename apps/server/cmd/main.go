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
	"github.com/lucjosin/qorv.in/internal/database/mariadb"
	"github.com/lucjosin/qorv.in/internal/slogx"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DatabaseConfig struct {
	Host         string `env:"HOST,required"`
	User         string `env:"USER,required"`
	Password     string `env:"PASSWORD,required"`
	Name         string `env:"NAME,required"`
	Port         int    `env:"PORT,required"`
	MaxOpenConns int    `env:"MAX_OPEN_CONN" envDefault:"5"`
	MaxIdleConns int    `env:"MAX_IDLE_CONN" envDefault:"2"`
}

type ServerConfig struct {
	Port           string `env:"PORT" envDefault:":4512"`
	URL            string `env:"URL"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`
}

type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Database DatabaseConfig `envPrefix:"DATABASE_"`
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

	db, err := mariadb.Open(ctx, mariadb.Config{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DatabaseName: cfg.Database.Name,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
	})
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error("closing database connection", "error", err)
		}
	}()

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
