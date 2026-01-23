package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/michael-duren/grind-75-cli/internal/config"
)

type SettingsModel struct {
	Inputs     []textinput.Model
	FocusIndex int
	Config     *config.Config
}

func NewSettingsModel() *SettingsModel {
	return &SettingsModel{
		Inputs: make([]textinput.Model, 0),
	}
}
