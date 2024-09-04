package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// migrationDirectory is the migrations' directory name
const migrationDirectory = "migrations"

// migrationFiles is populated when building the binary
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// migrateLatest runs the database migrations to the latest version.
func migrateLatest(db *sql.DB) error {
	// Build the migration migrationClient
	migrationClient, err := buildMigrationClient(db)
	if err != nil {
		return fmt.Errorf("migrate: failed to build migration migrationClient: %w", err)
	}

	// Defer closing the migrationClient
	defer func(m *migrate.Migrate) {
		sErr, tErr := m.Close()
		if sErr != nil {
			log.Printf("failed to close migration source: %v", sErr)
		}

		if tErr != nil {
			log.Printf("failed to close migration target: %v", tErr)
		}
	}(migrationClient)

	// Run the migrations to the latest version
	err = migrationClient.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: failed to run migration: %w", err)
	}

	// If err is migrate.ErrNoChange, then the migrationClient is already at the latest version, then log and return nil
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("DB Migration is up-to-date")
		return nil
	}

	log.Println("Successfully migrated DB to latest version")

	return nil
}

// buildMigrationClient creates a new migrate instance.
// source & target connection are required to be closed it in the calling function
func buildMigrationClient(db *sql.DB) (*migrate.Migrate, error) {
	// Read the source for the migrations. Our source is the SQL files in the migrations folder
	source, err := iofs.New(migrationFiles, migrationDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration source %w", err)
	}

	// Connect with the target i.e postgres DB
	target, err := pgx.WithInstance(db, new(pgx.Config))
	if err != nil {
		return nil, fmt.Errorf("failed to read migration target %w", err)
	}

	// Create a new instance of the migration using the defined source and target
	m, err := migrate.NewWithInstance("iofs", source, "postgres", target)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance %w", err)
	}

	return m, nil
}
