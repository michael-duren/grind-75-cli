package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitLogger(t *testing.T) {
	// Create temp directory for logs
	tempDir, err := os.MkdirTemp("", "g7clogs")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test initialization
	err = InitLogger(tempDir, true)
	if err != nil {
		t.Fatalf("InitLogger failed: %v", err)
	}

	// Verify log file creation
	expectedFile := filepath.Join(tempDir, "g7c-"+time.Now().Format("2006-01-02")+".log")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Log file not created at %s", expectedFile)
	}

	// Write something
	slog.Info("Test message")

	// Read file content and check length > 0
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if len(content) == 0 {
		t.Errorf("Log file empty after writing")
	}
}
