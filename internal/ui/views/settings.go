package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
)

func Settings(m *models.AppModel) string {
	if len(m.Settings.Inputs) == 0 {
		return "Loading settings..."
	}

	title := theme.ViewTitle.Render("Settings")

	// Hardcoded labels corresponding to input order in controller
	labels := []string{
		"Weeks for Grind",
		"Hours per Week",
		"SMTP Provider",
		"Email",
		"Password",
	}

	var content string
	for i, input := range m.Settings.Inputs {
		// Render label
		lbl := ""
		if i < len(labels) {
			lbl = theme.Label.Render(labels[i])
		}

		// Apply border style based on focus
		inputView := input.View()
		if i == m.Settings.FocusIndex {
			inputView = theme.InputFocused.Render(inputView)
		} else {
			inputView = theme.InputBlurred.Render(inputView)
		}

		content += fmt.Sprintf("\n%s\n%s\n", lbl, inputView)
		if i < len(m.Settings.Inputs)-1 {
			content += "\n"
		}
	}

	// Buttons
	buttons := []string{"Export Data", "Reset Settings", "Reset Data (WARNING)"}
	baseIdx := len(m.Settings.Inputs)

	for i, btn := range buttons {
		style := theme.BtnNormal
		prefix := "  "
		if m.Settings.FocusIndex == baseIdx+i {
			style = theme.BtnActive
			prefix = "> " // Already handled by style set string if using lipgloss properly, see below
			// Actually theme.BtnActive.SetString("> ") only affects what is returned if Render is called on empty string?
			// No, SetString is for list item bullets usually.
			// Let's stick to prefix string here for clarity.
		}
		content += fmt.Sprintf("\n\n%s%s", prefix, style.Render(btn))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		content,
		"\n\n(Esc to save & exit)",
	)
}
