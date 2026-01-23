package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/data/db"
	"github.com/michael-duren/grind-75-cli/internal/ui/controllers"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

type App struct {
	router *controllers.Router
}

func NewApp(services db.Service) *App {
	model := models.NewAppModel(services)

	return &App{
		router: controllers.NewRouter(model),
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := a.router.Route(msg)
	return a, cmd
}

func (a *App) View() string {
	return a.router.View()
}
