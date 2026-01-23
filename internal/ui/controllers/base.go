// internal/ui/controllers/controller.go
package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

type Controller interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
}

type Router struct {
	model       *models.AppModel
	controllers map[models.CurrentView]Controller
}

func NewRouter(m *models.AppModel) *Router {
	return &Router{
		model: m,
		controllers: map[models.CurrentView]Controller{
			models.HomePath:     NewHomeController(m),
			models.SettingsPath: NewSettingsController(m),
			models.HelpPath:     NewHelpController(m),
		},
	}
}

func (r *Router) Route(msg tea.Msg) tea.Cmd {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.model.Width = msg.Width
		r.model.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit
		case "s":
			r.model.CurrentView = models.SettingsPath
			return r.controllers[models.SettingsPath].Update(nil)
		case "?":
			r.model.CurrentView = models.HelpPath
			return r.controllers[models.HelpPath].Update(nil)
		}
	}

	// Route to current controller
	controller := r.controllers[r.model.CurrentView]
	if controller == nil {
		r.model.CurrentView = models.HomePath
		controller = r.controllers[models.HomePath]
	}

	return controller.Update(msg)
}

func (r *Router) View() string {
	controller := r.controllers[r.model.CurrentView]
	if controller == nil {
		return "Unknown view"
	}
	return controller.View()
}
