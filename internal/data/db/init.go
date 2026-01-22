package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string
	DB() *sql.DB
	// Conn() *pgx.Conn
	Queries() *dbgen.Queries

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

func InitServices(dbPath string) (Service, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set sqlite pragmas for better performance/safety
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Run migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Check if seeding is needed
	if err := checkAndSeed(db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return newService(db)
}

type service struct {
	db      *sql.DB
	queries *dbgen.Queries
}

func (s *service) Health() map[string]string {
	return map[string]string{
		"status": "ok",
	}
}

func (s *service) DB() *sql.DB {
	return s.db
}

func (s *service) Queries() *dbgen.Queries {
	return s.queries
}

func (s *service) Close() error {
	return s.db.Close()
}

func newService(db *sql.DB) (Service, error) {
	s := &service{
		db:      db,
		queries: dbgen.New(db),
	}

	return s, nil
}

func checkAndSeed(db *sql.DB) error {
	q := dbgen.New(db)
	ctx := context.Background()

	// Check if any problems exist
	probs, err := q.ListProblems(ctx)
	if err != nil {
		return err
	}

	if len(probs) == 0 {
		return Seed(ctx, db)
	}

	return nil
}
