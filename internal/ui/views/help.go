package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
)

func Help(m *models.AppModel) string {
	if !m.Help.TableInitialized {
		return "Loading..."
	}

	title := theme.ViewTitle.Render("Help & Keybindings")

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
