package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Help(m *models.AppModel) string {
	if !m.Help.TableInitialized {
		return "Loading..."
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(1, 0).
		Render("Help & Keybindings")

	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	tableStr := baseStyle.Render(m.Help.Table.View())

	footer := "\n(Press q or esc to return to Home)"

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		tableStr,
		footer,
	)
}
