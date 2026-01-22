package models

type HomeModel struct {
	// TODO: Add problem list state here
	Cursor int
}

func NewHomeModel() *HomeModel {
	return &HomeModel{
		Cursor: 0,
	}
}
