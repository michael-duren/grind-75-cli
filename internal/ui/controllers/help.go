package controllers

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Help(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	if !m.Help.TableInitialized {
		columns := []table.Column{
			{Title: "Key", Width: 20},
			{Title: "Action", Width: 40},
		}

		var rows []table.Row
		// Flatten help groups
		helpGroups := m.Help.Keys.FullHelp()
		for _, group := range helpGroups {
			for _, binding := range group {
				if len(binding.Keys()) > 0 {
					rows = append(rows, table.Row{
						binding.Help().Key,
						binding.Help().Desc,
					})
				}
			}
			// Add separator
			if len(group) > 0 {
				rows = append(rows, table.Row{"", ""})
			}
		}

		// Calc height
		height := max(m.Height-10, 5)

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(height),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m.Help.Table = t
		m.Help.TableInitialized = true
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.CurrentView = models.HomePath
			return m, nil
		}
	case tea.WindowSizeMsg:
		h := max(msg.Height-10, 5)
		m.Help.Table.SetHeight(h)
	}

	m.Help.Table, cmd = m.Help.Table.Update(msg)
	return m, cmd
}
