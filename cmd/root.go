package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/michael-duren/grind-75-cli/db"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "g7c",
	Short: "Grind 75 CLI Progress Tracker",
	Long:  `A terminal-based progress tracker for the Grind 75 coding challenge.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize DB
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		dbPath := filepath.Join(home, ".g7c", "g7c.db")
		database, err := db.InitDB(dbPath)
		if err != nil {
			fmt.Printf("Failed to initialize database: %v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		fmt.Println("Grind 75 CLI - Database initialized and ready.")
		// Start TUI here later...
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
