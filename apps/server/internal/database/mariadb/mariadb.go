package mariadb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lucjosin/qorv.in/internal/slogx"
	"github.com/lucjosin/qorv.in/migrations"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	driver "github.com/go-sql-driver/mysql"
)

type Database = sql.DB

const (
	sourceName = "iofs"
	driverName = "mysql"
)

// Config holds the configuration parameters for connecting to a MariaDB database.
type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	MaxOpenConns int
	MaxIdleConns int
}

func (c Config) DSN(allowMultiStatements bool) string {
	dsn := driver.Config{
		AllowNativePasswords: true,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", c.Host, c.Port),
		User:                 c.User,
		Passwd:               c.Password,
		DBName:               c.DatabaseName,
		ParseTime:            true,
		MultiStatements:      allowMultiStatements,
	}
	return dsn.FormatDSN()
}

// Open establishes a connection to the MariaDB database using the provided configuration.
func Open(ctx context.Context, config Config) (*Database, error) {
	log := slogx.FromCtx(ctx)
	log.Debug("connecting to database", "database", config.DatabaseName)

	// For security reasons, we disable multi-statements
	// by default to prevent SQL injection attacks.
	const allowMultiStatements = false
	db, err := sql.Open(driverName, config.DSN(allowMultiStatements))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = Migrate(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Info("connected to the database")
	return db, nil
}

// Migrate runs database migrations using the provided connection URL.
func Migrate(ctx context.Context, config Config) error {
	log := slogx.FromCtx(ctx)

	// For migrations, we need to allow multi-statements to execute
	// multiple SQL commands in a single migration file.
	const allowMultiStatements = true
	db, err := sql.Open(driverName, config.DSN(allowMultiStatements))
	if err != nil {
		return fmt.Errorf("opening database connection: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("creating migration driver: %v", err)
	}

	source, err := iofs.New(migrations.FS, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %v", err)
	}

	mg, err := migrate.NewWithInstance(sourceName, source, driverName, driver)
	if err != nil {
		return fmt.Errorf("configuring migrations: %v", err)
	}

	log.Debug("running database migrations")
	if err := mg.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to run migrations: %v", err)
		}
		log.Info("database schema already up to date, no migrations applied")
		return nil
	}

	log.Info("database migrations applied")
	return nil
}
