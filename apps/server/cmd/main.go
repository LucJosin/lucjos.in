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
	"github.com/lucjosin/qorv.in/internal/domain/domain"
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

	// domain
	domainRepo := domain.NewMariaDBRepository(db)
	domainService := domain.NewService(domainRepo)

	err = serverBootstrap(ctx, cfg, systemService, userService, domainService)
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
func serverBootstrap(ctx context.Context, cfg Config, systemService system.Service, userService user.Service, domainService domain.Service) error {
	log := slogx.FromCtx(ctx)
	log.Debug("bootstrapping server")

	var userEntity user.User
	systemEntity, err := systemService.Info(ctx)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return fmt.Errorf("loading system info: %w", err)
		}
		log.Info("system configuration not found, initiating first-time setup")

		// find or create the initial user
		userEntity, created, err := userService.FindOrCreateByUsername(ctx, user.User{
			FirstName:     cfg.App.Username,
			Username:      cfg.App.Username,
			Password:      cfg.App.Password, // TODO: hash the password
			Email:         cfg.App.Email,
			EmailVerified: true,
		})
		if err != nil {
			return fmt.Errorf("creating default user: %w", err)
		}
		if created {
			log.Info("system user created", "user", userEntity.Username)
		} else {
			log.Info("system user already exists, using existing account", "user", userEntity.Username)
		}

		// set this user as the immutable owner of the system
		systemEntity, err = systemService.Configure(ctx, userEntity.ID)
		if err != nil {
			return fmt.Errorf("configuring system: %w", err)
		}
		log.Info("system configured")
	}

	if (userEntity == user.User{}) {
		userEntity, err = userService.FindByID(ctx, systemEntity.OwnerUserID)
		if err != nil {
			return fmt.Errorf("loading system owner: %w", err)
		}
	}

	// after the initialization of a qorvin database, the owner/primary user cannot be changed.
	if userEntity.Username != cfg.App.Username {
		return errors.New("local configuration username does not match registered system owner")
	}

	_, created, err := domainService.FindOrCreateByDomain(ctx, domain.Domain{
		Domain: cfg.App.Domain,
	})
	if err != nil {
		return fmt.Errorf("configuring domain: %w", err)
	}
	if created {
		log.Info("system domain configured", "domain", cfg.App.Domain)
	} else {
		log.Info("system domain already exists, using existing one", "domain", cfg.App.Domain)
	}

	// TODO: setup workspace

	log.Info("server bootstrap complete")
	return nil
}
