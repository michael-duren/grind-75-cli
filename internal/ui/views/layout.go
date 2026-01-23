package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
)

func Layout(m *models.AppModel) string {
	// Header
	header := theme.
		AppHeader.
		Width(m.Width).
		Render("Grind 75 CLI")

	if m.Error != "" {
		// lipgloss.JoinVertical(lipgloss.Left,
		// 	header,
		// 	theme.ErrorStyle.Render(m.Error),
		// )

	}

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
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		"(q to quit, s for settings, ? for help)",
	)

	// Apply background to the whole view
	return lipgloss.NewStyle().
		// Background(theme.ColorBgLeetCode).
		Width(m.Width).
		Height(m.Height).
		Render(content)
}
