package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Home(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// In Home view, 'q' also quits
			return m, tea.Quit
		}
	}
	return m, nil
}
