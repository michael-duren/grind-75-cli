package controllers_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/ui/controllers"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func TestHome_Navigation(t *testing.T) {
	m := models.NewAppModel(nil)

	// Test 's' to go to settings
	newM, _ := controllers.Base(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

	if newM.CurrentView != models.SettingsPath {
		t.Errorf("Expected view to be Settings, got %s", newM.CurrentView)
	}
}

func TestHome_TableInput(t *testing.T) {
	m := models.NewAppModel(nil)

	// Mock some problems
	m.Home.Problems = []models.UserProblemWithRelations{
		{Title: "P1", URL: "http://p1.com"},
		{Title: "P2", URL: "http://p2.com"},
	}

	// Initialize table by calling Home with a dummy message
	// The first call should initialize the table
	newM, _ := controllers.Home(m, tea.WindowSizeMsg{Width: 100, Height: 50})

	if !newM.Home.TableInitialized {
		t.Error("Expected table to be initialized")
	}

	// Test 'j' to move down
	// Bubble Tea table handles 'j' internally, we just need to ensure the update cmd runs and state updates
	// But since Table is a model inside Home, we rely on Bubble Tea's update propagation.
	// Direct testing of bubble tea internal state (like table cursor) might be tricky without rendering or inspecting the inner model.
	// For now let's verify custom key handling like 'l' for column selection.

	initialCol := newM.Home.SelectedCol
	newM, _ = controllers.Home(newM, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if newM.Home.SelectedCol == initialCol {
		t.Error("Expected SelectedCol to change on 'l'")
	}
}

func TestSettings_Navigation(t *testing.T) {
	m := models.NewAppModel(nil)
	m.CurrentView = models.SettingsPath

	// Test 'esc' to go back home
	newM, _ := controllers.Settings(m, tea.KeyMsg{Type: tea.KeyEscape})

	if newM.CurrentView != models.HomePath {
		t.Errorf("Expected view to be Home, got %s", newM.CurrentView)
	}
}
