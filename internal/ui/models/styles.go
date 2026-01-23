package models

import "github.com/charmbracelet/lipgloss"

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	NoStyle      = lipgloss.NewStyle()
	CursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)
