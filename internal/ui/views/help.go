package views

import (
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Help(m *models.AppModel) string {
	helpText := `
		Grind 75 CLI Help

		Navigation:
		- From Home View:
		  - Press 's' to go to Settings View
		  - Press '?' to go to Help View
		  - Press 'q' to quit the application
		- Page Navigation:
		  - Use 'up' and 'down' arrow keys to scroll through content if applicable
		  - Press 'space' to select options

		- From Settings View:
		  - Press 'esc' to return to Home View

		- From Help View:
		  - Press 'esc' to return to Home View

		General Commands:
		- 'ctrl+c': Quit the application from any view

		Enjoy using Grind 75 CLI!
		`
	return helpText
}
