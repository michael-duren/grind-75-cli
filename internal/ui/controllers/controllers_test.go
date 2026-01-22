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
	newM, _ := controllers.Home(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

	if newM.CurrentView != models.SettingsPath {
		t.Errorf("Expected view to be Settings, got %s", newM.CurrentView)
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
