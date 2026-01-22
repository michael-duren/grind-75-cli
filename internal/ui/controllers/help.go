package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Help(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.CurrentView = models.HomePath
			return m, nil
		}
	}
	return m, nil
}
