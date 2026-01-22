package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/data/db"
	"github.com/michael-duren/grind-75-cli/internal/logging"
	"github.com/michael-duren/grind-75-cli/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "g7c",
	Short: "Grind 75 CLI Progress Tracker",
	Long:  `A terminal-based progress tracker for the Grind 75 coding challenge.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Logger
		err := logging.InitLogger("", false)
		if err != nil {
			fmt.Printf("Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		// Initialize DB
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("Failed to get home directory", "err", err)
			os.Exit(1)
		}
		dbPath := filepath.Join(home, ".g7c", "grind75.db")

		database, err := db.InitDB(dbPath)
		if err != nil {
			slog.Error("Failed to initialize database", "err", err)
			os.Exit(1)
		}
		defer database.Close()

		slog.Info("Grind 75 CLI Initialized")

		// Start TUI
		p := tea.NewProgram(ui.NewApp(database), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	execute()
}
