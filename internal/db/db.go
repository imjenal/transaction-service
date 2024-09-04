package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB is a wrapper for database connection
type DB struct {
	Conn *pgxpool.Pool
}

// Config is a configuration for database connection
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	ForceTLS bool
	Migrate  bool
	AppName  string
	SeedDB   bool
}

// Link sqlc with go generate, now we need to just run go generate to generate models and functions for DB
//go:generate sqlc generate

//go:generate mockgen -source=models/querier.go -destination=models/mock/querier.go -package=mock

var (
	once   sync.Once
	dbConn *DB
)

// GetConnection creates a new connection if not available or else returns the current connection
func GetConnection(ctx context.Context, cfg *Config) (*DB, error) {
	var (
		conn *DB
		err  error
	)

	// Use a sync.Once to make sure we only create the connection once
	once.Do(func() {
		// Make the database connection, panic if we can't connect
		conn, err = connect(ctx, cfg)
		if err != nil {
			log.Printf("failed to connect database: %v", err)
			return
		}

		// Set the global connection, it will be used by the rest of the application. Future calls to GetConnection will return this connection
		dbConn = conn

		// If we don't need to migrate the database, return
		if !cfg.Migrate {
			return
		}

		// Create a new sql.DB connection, it will be used by the migration library to run migrations.
		//We don't want to use the pgx connection because the migration closes the connection after running the migrations,
		// and we don't want to close the pgx connection as it will be used by the rest of the application
		sqlDB, err := sql.Open("pgx", cfg.ConnString())
		if err != nil {
			log.Printf("failed to connect database: %v", err)
			return
		}

		// Migrate the database to the latest version
		if err = migrateLatest(sqlDB); err != nil {
			log.Printf("failed to migrate database: %v", err)
			return
		}
	})

	// Seed the database if needed(only in development, when SeedDB is true)
	if cfg.SeedDB {
		log.Printf("Seeding database...")

		files, err := seedDB(ctx, conn)
		if err != nil {
			log.Printf("failed to seed database")
			return dbConn, err
		}

		log.Printf("Database seeded successfully: %v", files)
	}

	return dbConn, err
}

// connect creates a new DB connection
func connect(ctx context.Context, cfg *Config) (*DB, error) {
	// Build the pgx connection config
	pgxConfig, err := getPgxConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Create the pgx connection pool
	conn, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &DB{
		Conn: conn,
	}, nil
}

// getPgxConfig builds and returns the pgx connection config
func getPgxConfig(cfg *Config) (*pgxpool.Config, error) {
	// Build the pgx connection config
	connConfig, err := pgxpool.ParseConfig(cfg.ConnString())
	if err != nil {
		return nil, fmt.Errorf("getPgxConfig: %w", err)
	}

	// Connection pool configuration
	connConfig.MaxConns = 5
	connConfig.MinConns = 2

	return connConfig, nil
}

// ConnString returns a Postgres connection URL e.g.,
// postgres://username:password@localhost:5432/database_name?sslmode=disable
func (cfg *Config) ConnString() string {
	// Set default ssl mode to prefer
	sslMode := "prefer"

	// If force tls is set, set ssl mode to require
	if cfg.ForceTLS {
		sslMode = "require"
	}

	// Build the postgres connection URL
	dbURL := url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		User:     url.UserPassword(cfg.User, cfg.Password),
		Path:     cfg.Name,
		RawQuery: "sslmode=" + sslMode,
	}

	// If the app name is set, add it to the connection URL
	if cfg.AppName != "" {
		dbURL.RawQuery += "&application_name=" + cfg.AppName
	}

	return dbURL.String()
}
