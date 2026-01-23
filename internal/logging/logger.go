package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

// InitLogger initializes the global Logger to write to both stdout and a file.
// If logDir is empty, it defaults to ~/.g7c/logs.
func InitLogger(logDir string, debug bool) error {
	// Determine log file path
	if logDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home dir: %w", err)
		}
		logDir = filepath.Join(home, ".g7c", "logs")
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log dir: %w", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("g7c-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	charmHandler := log.New(io.Writer(f))
	charmHandler.SetReportCaller(true)

	if debug {
		charmHandler.SetLevel(log.DebugLevel)
	} else {
		charmHandler.SetLevel(log.InfoLevel)
	}

	logger := slog.New(charmHandler)
	slog.SetDefault(logger)

	slog.Info("Logger initialized", "log_file", logFile, "debug_mode", debug)
	return nil
}
