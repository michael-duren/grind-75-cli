package models

import "github.com/charmbracelet/bubbles/table"

type HelpModel struct {
	Table            table.Model
	TableInitialized bool
	Keys             KeyMap
}

func NewHelpModel() *HelpModel {
	return &HelpModel{
		Keys: Keys,
	}
}
