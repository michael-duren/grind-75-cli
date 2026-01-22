package models

import (
	"database/sql"
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
	DB *sql.DB

	*WindowDimensions
	BodyDimensions *WindowDimensions

	CurrentView CurrentView

	// Sub-Models
	Home     *HomeModel
	Settings *SettingsModel
}

func NewAppModel(db *sql.DB) *AppModel {
	return &AppModel{
		DB:               db,
		WindowDimensions: &WindowDimensions{},
		BodyDimensions:   &WindowDimensions{},
		CurrentView:      HomePath,

		Home:     NewHomeModel(),
		Settings: NewSettingsModel(),
	}
}
