package views

import tea "github.com/charmbracelet/bubbletea"

// View defines the interface that all sub-views must implement.
type View interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}
