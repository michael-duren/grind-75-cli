package db_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/michael-duren/grind-75-cli/db"
	gen "github.com/michael-duren/grind-75-cli/db/gen"
)

func TestInitDB(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "g7ctest")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")

	// InitDB should create DB, run migrations, and seed
	database, err := db.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer database.Close()

	// Verify tables exist (by querying problems)
	q := gen.New(database)
	probs, err := q.ListProblems(context.Background())
	if err != nil {
		t.Fatalf("ListProblems failed: %v", err)
	}

	if len(probs) == 0 {
		t.Errorf("Expected problems to be seeded, got 0")
	}

	// Verify idempotency (run InitDB again)
	database2, err := db.InitDB(dbPath)
	if err != nil {
		t.Fatalf("Second InitDB failed: %v", err)
	}
	database2.Close()

	// Check if data is not duplicated (should be same count)
	probs2, err := q.ListProblems(context.Background())
	if err != nil {
		t.Fatalf("ListProblems 2 failed: %v", err)
	}
	if len(probs) != len(probs2) {
		t.Errorf("Expected problem count to match %d, got %d", len(probs), len(probs2))
	}
}
