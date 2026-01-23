package models

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
)

type UserProblemWithRelations struct {
	ProblemID       int64
	Slug            string
	Title           string
	URL             string
	Duration        int64
	DifficultyID    string
	DifficultyName  string
	Status          string
	LastAttemptedAt *time.Time
	Attempts        int64
	Topics          []dbgen.Topic
	Reviews         []dbgen.Review
}

type HomeModel struct {
	Problems         []UserProblemWithRelations
	Table            table.Model
	TableInitialized bool
	SelectedCol      int

	// Details / Editing
	NotesInput      textarea.Model
	Editing         bool
	ActiveProblemID int64
}

func NewHomeModel() *HomeModel {
	ta := textarea.New()
	ta.Placeholder = "Add notes here..."
	return &HomeModel{
		Problems:   []UserProblemWithRelations{},
		NotesInput: ta,
	}
}
