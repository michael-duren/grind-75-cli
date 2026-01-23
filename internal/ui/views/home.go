package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
)

func Home(m *models.AppModel) string {
	if !m.Home.TableInitialized {
		return "Loading..."
	}

	// Calculate widths
	totalWidth := m.Width
	// Estimate table width (sum of cols + borders + padding)
	// Cols: 4+25+8+12+6+4 = 59. Plus borders/padding ~10. Total ~70.
	// Layout logic
	// Minimum width required for the table to look good (sum of columns ~60 + padding)
	minTableWidth := 72

	tableWidth := totalWidth
	detailsWidth := totalWidth
	stackVertical := true

	if totalWidth >= 110 {
		stackVertical = false
		// Give table roughly 70% of width, but ensure it's at least minTableWidth
		tableWidth = max(totalWidth*70/100, minTableWidth)
		// Details gets the rest
		detailsWidth = totalWidth - tableWidth - 2
	} else {
		// Stacked mode
		tableWidth = totalWidth - 4
		detailsWidth = totalWidth - 4
	}

	// Render Table (Force width on style just in case)
	tableView := theme.Base.Width(tableWidth).Render(m.Home.Table.View())

	// Render Details or Edit
	var detailsView string

	idx := m.Home.Table.Cursor()
	if idx >= 0 && idx < len(m.Home.Problems) {
		p := m.Home.Problems[idx]

		if m.Home.Editing {
			title := theme.ViewTitle.Render("Edit Notes")
			m.Home.NotesInput.SetWidth(detailsWidth - 2)
			m.Home.NotesInput.SetHeight(m.Height - 20) // Approx
			input := m.Home.NotesInput.View()
			detailsView = lipgloss.JoinVertical(lipgloss.Left, title, input, "\n(Esc to cancel, Ctrl+S to save)")
		} else {
			// View Mode
			title := theme.ViewTitle.Render(p.Title)
			info := fmt.Sprintf("Difficulty: %s\nTime: %d mins\nTopics: ", p.DifficultyName, p.Duration)
			if len(p.Topics) > 0 {
				for i, t := range p.Topics {
					if i > 0 {
						info += ", "
					}
					info += t.Name
				}
			}
			info += "\n\n"

			notes := "No notes."
			// Find latest review notes
			if len(p.Reviews) > 0 {
				// Assuming reviews are ordered desc by date
				if p.Reviews[0].Notes.Valid {
					notes = p.Reviews[0].Notes.String
				}
			}

			notesHeader := theme.Label.Render("Notes:")
			notesBody := theme.Base.Render(notes)

			detailsView = lipgloss.JoinVertical(lipgloss.Left,
				title,
				theme.Base.Render(info),
				notesHeader,
				notesBody,
				"\n(Enter to edit notes, 'o' to open)",
			)
		}
	} else {
		detailsView = "Select a problem..."
	}

	helpView := m.Home.CommandBar.View()

	tableWithHelp := lipgloss.JoinVertical(
		lipgloss.Left,
		tableView,
		"",
		helpView,
	)

	// Container styles
	detailsStyle := lipgloss.NewStyle().
		Width(detailsWidth).
		Height(m.Height-20).
		Padding(1).
		Margin(1, 0).
		BorderForeground(theme.ColorBrand)

	if stackVertical {
		detailsStyle = detailsStyle.Height(10)
		return lipgloss.JoinVertical(lipgloss.Left,
			tableWithHelp,
			detailsStyle.Render(detailsView),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		tableWithHelp,
		detailsStyle.Render(detailsView),
	)
}
