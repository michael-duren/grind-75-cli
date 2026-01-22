package models

import (
	"github.com/michael-duren/grind-75-cli/internal/data/db"
)

type WindowDimensions struct {
	Width  int
	Height int
}

type CurrentView string

const (
	HomePath     CurrentView = "/home"
	SettingsPath CurrentView = "/settings"
	HelpPath     CurrentView = "/help"
)

type AppModel struct {
	Services db.Service

	*WindowDimensions
	BodyDimensions *WindowDimensions

	CurrentView CurrentView

	// Sub-Models
	Home     *HomeModel
	Settings *SettingsModel
}

func NewAppModel(services db.Service) *AppModel {
	return &AppModel{
		Services:         services,
		WindowDimensions: &WindowDimensions{},
		BodyDimensions:   &WindowDimensions{},
		CurrentView:      HomePath,

		Home:     NewHomeModel(),
		Settings: NewSettingsModel(),
	}
}
