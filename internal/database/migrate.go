package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

// MigrationsFS is set by main.go to embed migrations from root directory
var MigrationsFS embed.FS

// Migrate runs all pending migrations
func Migrate(db *sql.DB) error {
	// Create migrations tracking table if not exists
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := fs.ReadDir(MigrationsFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by name (they're numbered, so alphabetical = chronological)
	var migrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	// Get already-run migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, filename := range migrations {
		if applied[filename] {
			continue
		}

		log.Printf("Running migration: %s", filename)

		content, err := fs.ReadFile(MigrationsFS, filename)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		if err := runMigration(db, filename, string(content)); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", filename, err)
		}

		log.Printf("Completed migration: %s", filename)
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT filename FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		applied[filename] = true
	}
	return applied, rows.Err()
}

func runMigration(db *sql.DB, filename, content string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Run the migration SQL
	if _, err := tx.Exec(content); err != nil {
		return err
	}

	// Record that we ran it
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (filename) VALUES ($1)",
		filename,
	); err != nil {
		return err
	}

	return tx.Commit()
}
