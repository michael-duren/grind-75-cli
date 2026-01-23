package controllers

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/config"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
)

const (
	inputWeeks = iota
	inputHours
	inputProvider
	inputEmail
	inputPassword
	// Virtual buttons
	btnExport
	btnResetSettings
	btnResetData
)

func Settings(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	// Initialize inputs if empty
	if len(m.Settings.Inputs) == 0 {
		cfg, err := config.LoadConfig()
		if err != nil {
			slog.Error("Failed to load config", "error", err)
			cfg = config.GetDefault()
		}
		m.Settings.Config = cfg
		initInputs(m)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			saveInputsToConfig(m)
			m.CurrentView = models.HomePath
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// If enter on a button, execute action
			if s == "enter" && m.Settings.FocusIndex >= btnExport {
				return executeAction(m)
			}

			// Navigation
			if s == "up" || s == "shift+tab" {
				m.Settings.FocusIndex--
			} else {
				m.Settings.FocusIndex++
			}

			// Cycle
			if m.Settings.FocusIndex > btnResetData {
				m.Settings.FocusIndex = 0
			} else if m.Settings.FocusIndex < 0 {
				m.Settings.FocusIndex = btnResetData
			}

			// Update focus
			cmds := make([]tea.Cmd, len(m.Settings.Inputs))
			for i := 0; i <= len(m.Settings.Inputs)-1; i++ {
				if i == m.Settings.FocusIndex {
					// Set focused
					cmds[i] = m.Settings.Inputs[i].Focus()
					// We need to apply Text styles to the input model
					// But wait, the VIEW renders the labels. The INPUT view renders the box.
					// textinput.Model has TextStyle (text) and CursorStyle.
					// It doesn't natively have a "Border Style" property that wraps the *entire* input unless we wrap it in View.
					// However, standard textinput view is just the text area.
					// If we want a border, we usually wrap the View() output.
					// BUT the `newInput` can be configured.

					// Let's rely on the View wrapping or just color the text for now as per previous logic.
					// Actually, standard Bubbles textinput doesn't have a border by default.
					// To get the "LeetCode Input" look (Bordered), we should apply a border style in the VIEW, not here.
					// Here we just toggle focus state.
					m.Settings.Inputs[i].TextStyle = theme.Base
					m.Settings.Inputs[i].Cursor.Style = theme.InputFocused
				} else {
					m.Settings.Inputs[i].Blur()
					m.Settings.Inputs[i].TextStyle = theme.Label // Blurred text gray
					m.Settings.Inputs[i].Cursor.Style = theme.InputBlurred
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	// Update inputs (only the focused one needs update really, but batching is fine)
	cmdBatch := make([]tea.Cmd, len(m.Settings.Inputs))
	for i := range m.Settings.Inputs {
		m.Settings.Inputs[i], cmdBatch[i] = m.Settings.Inputs[i].Update(msg)
	}

	return m, tea.Batch(cmdBatch...)
}

func initInputs(m *models.AppModel) {
	cfg := m.Settings.Config
	inputs := make([]textinput.Model, 5)

	inputs[inputWeeks] = newInput("Weeks for Grind", strconv.Itoa(cfg.GrindPlan.Weeks), "4")
	inputs[inputHours] = newInput("Hours per Week", strconv.Itoa(cfg.GrindPlan.HoursPerWeek), "8")
	inputs[inputProvider] = newInput("SMTP Provider", cfg.SMTP.Provider, "gmail")
	inputs[inputEmail] = newInput("Email", cfg.SMTP.Email, "user@example.com")
	inputs[inputPassword] = newInput("Password", cfg.SMTP.Password, "secret")
	inputs[inputPassword].EchoMode = textinput.EchoPassword

	m.Settings.Inputs = inputs
	m.Settings.FocusIndex = 0
	m.Settings.Inputs[0].Focus()
}

func newInput(placeholder, value, example string) textinput.Model {
	t := textinput.New()
	t.Placeholder = placeholder
	t.SetValue(value)
	t.CharLimit = 156
	t.Width = 30
	return t
}

func saveInputsToConfig(m *models.AppModel) {
	// Parse back
	if w, err := strconv.Atoi(m.Settings.Inputs[inputWeeks].Value()); err == nil {
		m.Settings.Config.GrindPlan.Weeks = w
	}
	if h, err := strconv.Atoi(m.Settings.Inputs[inputHours].Value()); err == nil {
		m.Settings.Config.GrindPlan.HoursPerWeek = h
	}
	m.Settings.Config.SMTP.Provider = m.Settings.Inputs[inputProvider].Value()
	m.Settings.Config.SMTP.Email = m.Settings.Inputs[inputEmail].Value()
	m.Settings.Config.SMTP.Password = m.Settings.Inputs[inputPassword].Value()

	if err := m.Settings.Config.SaveConfig(); err != nil {
		slog.Error("Failed to save config", "error", err)
	}
}

func executeAction(m *models.AppModel) (*models.AppModel, tea.Cmd) {
	switch m.Settings.FocusIndex {
	case btnExport:
		problems, err := m.Services.Queries().ListUserProblems(context.Background())
		if err == nil {
			data, _ := json.MarshalIndent(problems, "", "  ")
			_ = os.WriteFile("grind75_export.json", data, 0644)
		}
	case btnResetSettings:
		m.Settings.Config = config.GetDefault()
		_ = m.Settings.Config.SaveConfig()
		initInputs(m) // Re-init inputs with defaults
	case btnResetData:
		// Not implemented
	}
	return m, nil
}
