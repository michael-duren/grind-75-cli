package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	dbgen "github.com/michael-duren/grind-75-cli/internal/data/db/gen"
)

func TestSeed(t *testing.T) {
	// Create an in-memory database
	// We need to apply the schema first.
	// We can read schema.sql and execute it.

	dbConn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer dbConn.Close()

	// Load schema
	schema, err := os.ReadFile("./schema.sql")
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}

	if _, err := dbConn.Exec(string(schema)); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}

	ctx := context.Background()

	// Run Seed
	if err := Seed(ctx, dbConn); err != nil {
		t.Fatalf("Seed failed: %v", err)
	}

	queries := dbgen.New(dbConn)

	// Verify problems count (should be 169 based on default json usually, or around 75+extended)
	// problems.json has 828 lines, likely ~75-80 items?
	// Let's just check > 0

	probs, err := queries.ListProblems(ctx)
	if err != nil {
		t.Fatalf("ListProblems failed: %v", err)
	}

	if len(probs) == 0 {
		t.Errorf("Expected problems to be seeded, got 0")
	}

	// Verify a specific one, e.g., Two Sum
	twoSum, err := queries.GetProblemBySlug(ctx, "two-sum")
	if err != nil {
		t.Fatalf("GetProblemBySlug failed: %v", err)
	}
	if twoSum.Title != "Two Sum" {
		t.Errorf("Expected Two Sum, got %s", twoSum.Title)
	}

	// Verify user progress initialized
	prog, err := queries.GetUserProgress(ctx, twoSum.ID)
	if err != nil {
		t.Fatalf("GetUserProgress failed: %v", err)
	}
	if prog.Status != "New" {
		t.Errorf("Expected status New, got %s", prog.Status)
	}
}
