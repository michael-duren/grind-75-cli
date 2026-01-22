package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Home(m *models.AppModel) string {
	if !m.Home.TableInitialized {
		return "Loading..."
	}

	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	// Render table
	tableStr := baseStyle.Render(m.Home.Table.View())

	// Highlight the header of the selected column to indicate "edit mode" for that column (future feature)
	// For now just show which column is selected for potential action
	colHeaders := []string{"Status", "Problem", "Difficulty", "Topic", "Time", "Attempts"}
	headerView := ""
	for i, h := range colHeaders {
		style := lipgloss.NewStyle().Padding(0, 1)
		if i == m.Home.SelectedCol {
			style = style.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
		}
		headerView += style.Render(h)
	}

	helpText := fmt.Sprintf("\nSelected Column: %s (Space to interact)\n(q to quit, s for Settings, Use Arrows/jkhl to navigate)", colHeaders[m.Home.SelectedCol])

	return lipgloss.JoinVertical(lipgloss.Left,
		tableStr,
		helpText,
	)
}
