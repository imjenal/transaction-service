package db

import (
	"context"
	"embed"
	"fmt"
	"os"
	"sort"

	"github.com/jackc/pgx/v4"
)

const seedsDir = "seeds"

// seeds is populated when building the binary
//
//go:embed seeds/*.sql
var seeds embed.FS

// seedDB seeds the database with the queries in the seeds directory
func seedDB(ctx context.Context, db *DB) ([]string, error) {
	dbTxn, err := db.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("seedDB: failed to begin DB transaction: %w", err)
	}

	defer func(dbTxn pgx.Tx, ctx context.Context) {
		_ = dbTxn.Rollback(ctx)
	}(dbTxn, ctx)

	// Setup seed table only if it doesn't exist
	err = setupSeedTable(ctx, dbTxn)
	if err != nil {
		return nil, fmt.Errorf("seedDB: failed to setup seed table: %w", err)
	}

	// Get all files in the seeds directory
	files, err := seeds.ReadDir(seedsDir)
	if err != nil {
		return nil, fmt.Errorf("seedDB: failed to read seeds directory: %w", err)
	}

	// Some seed files need to run in order to prevent constraint errors
	// Even though fs.ReadDir sends us an ascending list by default,
	// we sort it manually to be safe
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// seeded variables keeps track of all the files that were seeded
	seeded := make([]string, 0, len(files))
	// Loop through all the files in the seeds directory
	for _, file := range files {
		// If the file is a directory, skip it
		if file.IsDir() {
			continue
		}

		// Execute the seed query in the file
		executed, err := executeSeedQuery(ctx, dbTxn, file.Name())
		if err != nil {
			return nil, fmt.Errorf("seedDB: failed to execute seed query for %s: %w", file.Name(), err)
		}

		// If the query was executed, add it to the seeded slice,
		// Also add it to the batch to mark it as executed in DB
		if executed {
			seeded = append(seeded, file.Name())
		}
	}

	// Mark all the seeded files as executed in DB
	if err = markExecuted(ctx, dbTxn, seeded); err != nil {
		return nil, fmt.Errorf("seedDB: failed to mark seeds as executed: %w", err)
	}

	if err = dbTxn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("seedDB: failed to commit DB transaction: %w", err)
	}

	return seeded, nil
}

// setupSeedTable creates the seeds table if it doesn't exist
func setupSeedTable(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS seeds (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return fmt.Errorf("setupSeedTable: failed to create seed table: %w", err)
	}

	return nil
}

// isAlreadyExecuted checks if the seed query in the file was already executed
func isAlreadyExecuted(ctx context.Context, tx pgx.Tx, fileName string) (bool, error) {
	var count int

	err := tx.QueryRow(ctx, "SELECT COUNT(id) FROM seeds WHERE name = $1", fileName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("isAlreadyExecuted: failed to query seed table: %w", err)
	}

	return count > 0, nil
}

// executeSeedQuery executes the seed query in the file if it was not executed before
func executeSeedQuery(ctx context.Context, tx pgx.Tx, fileName string) (bool, error) {
	// Check if the seed query was already executed
	executed, err := isAlreadyExecuted(ctx, tx, fileName)
	if err != nil {
		return false, fmt.Errorf("executeSeedQuery: failed to check if seed was already executed: %w", err)
	}

	// If the seed query was already executed, return false
	if executed {
		return false, nil
	}

	// Build the path to the seed file
	file := fmt.Sprintf("%s%s%s", seedsDir, string(os.PathSeparator), fileName)

	// Read the seed query from the file
	contents, err := seeds.ReadFile(file)
	if err != nil {
		return false, fmt.Errorf("executeSeedQuery: failed to read seed file: %w", err)
	}

	// Execute the seed query
	_, err = tx.Exec(ctx, string(contents))
	if err != nil {
		return false, fmt.Errorf("executeSeedQuery: failed to seed DB: %w", err)
	}

	// Return true to indicate that the seed query was executed
	return true, nil
}

func markExecuted(ctx context.Context, tx pgx.Tx, files []string) error {
	query := "INSERT INTO seeds (name) VALUES ($1)"

	for _, file := range files {
		_, err := tx.Exec(ctx, query, file)
		if err != nil {
			return fmt.Errorf("markExecuted: failed to mark seed as executed: %w", err)
		}
	}

	return nil
}
