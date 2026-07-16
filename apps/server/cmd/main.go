package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucjosin/qorv.in/internal/api"
	"github.com/lucjosin/qorv.in/internal/database/mariadb"
	"github.com/lucjosin/qorv.in/internal/domain/system"
	"github.com/lucjosin/qorv.in/internal/domain/user"
	"github.com/lucjosin/qorv.in/internal/errs"
	"github.com/lucjosin/qorv.in/internal/slogx"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppConfig struct {
	Domain        string `env:"DOMAIN,required"`
	Username      string `env:"USERNAME,required"`
	Password      string `env:"PASSWORD,required"`
	Email         string `env:"EMAIL,required"`
	WorkspaceName string `env:"WORKSPACE_NAME" envDefault:"Main"`
}

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
	App      AppConfig      `envPrefix:"APP_"`
	Database DatabaseConfig `envPrefix:"DATABASE_"`
	Server   ServerConfig   `envPrefix:"SERVER_"`
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

	// system
	systemRepo := system.NewMariaDBRepository(db)
	systemService := system.NewService(systemRepo)

	// user
	userRepo := user.NewMariaDBRepository(db)
	userService := user.NewService(userRepo)

	err = serverBootstrap(ctx, cfg, systemService, userService)
	if err != nil {
		log.Panic(err)
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

// serverBootstrap initializes the server.
func serverBootstrap(ctx context.Context, cfg Config, systemService system.Service, userService user.Service) error {
	log := slogx.FromCtx(ctx)
	log.Debug("bootstrapping server")

	systemEntity, err := systemService.Info(ctx)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return fmt.Errorf("loading system info: %w", err)
		}
		log.Info("system configuration not found, initiating first-time setup")

		// find or create the initial bootstrap user
		userEntity, created, err := userService.FindOrCreateByUsername(ctx, user.User{
			FirstName: cfg.App.Username,
			Username:  cfg.App.Username,
			Password:  cfg.App.Password,
			Email:     cfg.App.Email,
		})
		if err != nil {
			return fmt.Errorf("creating default bootstrap user: %w", err)
		}
		if created {
			log.Info("default bootstrap user created", "user", userEntity.Username)
		} else {
			log.Info("default bootstrap user already exists, using existing account", "user", userEntity.Username)
		}

		// TODO: setup workspace and domain

		// set this user as the immutable owner of the system
		err = systemService.Configure(ctx, userEntity.ID)
		if err != nil {
			return fmt.Errorf("configuring system: %w", err)
		}
		log.Info("system successfully initialized")

		return nil
	}

	userEntity, err := userService.FindByID(ctx, systemEntity.OwnerUserID)
	if err != nil {
		return fmt.Errorf("loading system owner: %w", err)
	}

	if userEntity.Username != cfg.App.Username {
		return fmt.Errorf("local configuration username does not match registered system owner")
	}

	log.Info("server bootstrap complete")
	return nil
}
