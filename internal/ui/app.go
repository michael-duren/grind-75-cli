package ui

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/ui/controllers"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/views"
)

type Model struct {
	AppModel *models.AppModel
}

func NewApp(db *sql.DB) Model {
	return Model{
		AppModel: models.NewAppModel(db),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := controllers.Base(m.AppModel, msg)
	m.AppModel = newModel
	return m, cmd
}

func (m Model) View() string {
	return views.Layout(m.AppModel)
}
