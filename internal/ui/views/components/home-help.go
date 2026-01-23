// internal/ui/components/command_bar.go
package components

import (
	"github.com/charmbracelet/lipgloss"
)

type CommandBar struct {
	commands []Command
	style    lipgloss.Style
}

type Command struct {
	Key  string
	Desc string
}

func NewCommandBar(commands []Command) CommandBar {
	return CommandBar{
		commands: commands,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1),
	}
}

func (c CommandBar) View() string {
	var parts []string

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246"))

	for _, cmd := range c.commands {
		part := keyStyle.Render(cmd.Key) + " " + descStyle.Render(cmd.Desc)
		parts = append(parts, part)
	}

	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("238")).
		Render(" • ")

	return c.style.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, parts[0]),
		separator,
		lipgloss.JoinHorizontal(lipgloss.Left, parts[1:]...),
	))
}
