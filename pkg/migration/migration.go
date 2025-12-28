package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Initialize creates the migrations table if it doesn't exist
func (m *Migrator) Initialize() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

// LoadMigrations loads all migration files from the migrations directory
func (m *Migrator) LoadMigrations() ([]Migration, error) {
	files, err := ioutil.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse version from filename (e.g., "001_create_users.sql")
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Warning: invalid migration filename %s", file.Name())
			continue
		}

		name := strings.TrimSuffix(parts[1], ".sql")

		// Read SQL content
		content, err := ioutil.ReadFile(filepath.Join(m.migrationsDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(content),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// GetAppliedMigrations returns list of applied migrations
func (m *Migrator) GetAppliedMigrations() (map[int]Migration, error) {
	query := `SELECT version, name, applied_at FROM schema_migrations ORDER BY version`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]Migration)
	for rows.Next() {
		var migration Migration
		err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		applied[migration.Version] = migration
	}

	return applied, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.Initialize(); err != nil {
		return err
	}

	migrations, err := m.LoadMigrations()
	if err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if _, exists := applied[migration.Version]; exists {
			log.Printf("Migration %03d_%s already applied, skipping", migration.Version, migration.Name)
			continue
		}

		log.Printf("Applying migration %03d_%s", migration.Version, migration.Name)

		// Execute migration in transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute the migration SQL
		_, err = tx.Exec(migration.UpSQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %03d_%s: %w", migration.Version, migration.Name, err)
		}

		// Record the migration as applied
		_, err = tx.Exec(
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			migration.Version, migration.Name,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %03d_%s: %w", migration.Version, migration.Name, err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %03d_%s: %w", migration.Version, migration.Name, err)
		}

		log.Printf("Successfully applied migration %03d_%s", migration.Version, migration.Name)
	}

	return nil
}

// Status shows migration status
func (m *Migrator) Status() error {
	if err := m.Initialize(); err != nil {
		return err
	}

	migrations, err := m.LoadMigrations()
	if err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, migration := range migrations {
		status := "PENDING"
		appliedAt := ""

		if appliedMigration, exists := applied[migration.Version]; exists {
			status = "APPLIED"
			if appliedMigration.AppliedAt != nil {
				appliedAt = appliedMigration.AppliedAt.Format("2006-01-02 15:04:05")
			}
		}

		fmt.Printf("%03d %-30s %s %s\n", migration.Version, migration.Name, status, appliedAt)
	}

	return nil
}
