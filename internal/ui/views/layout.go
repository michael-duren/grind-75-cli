package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Layout(m *models.AppModel) string {
	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(1).
		Render("Grind 75 CLI")

	// Body selection
	var body string
	switch m.CurrentView {
	case models.HomePath:
		body = Home(m)
	case models.SettingsPath:
		body = Settings(m)
	case models.HelpPath:
		body = Help(m)
	default:
		body = "Unknown View"
	}

	// Wrapper
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		"",
		"(q to quit, s for settings, ? for help)",
	)
}
