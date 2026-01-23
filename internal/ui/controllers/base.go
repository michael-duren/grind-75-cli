package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Base(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "s":
			m.CurrentView = models.SettingsPath
			return Settings(m, nil)
		case "?":
			m.CurrentView = models.HelpPath
			// Call Help with nil msg to ensure initialization runs immediately
			return Help(m, nil)
		}
	}

	switch m.CurrentView {
	case models.HomePath:
		return Home(m, msg)
	case models.SettingsPath:
		return Settings(m, msg)
	case models.HelpPath:
		return Help(m, msg)
	}

	return m, nil
}
