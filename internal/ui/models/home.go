package models

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
	"github.com/michael-duren/grind-75-cli/internal/ui/views/components"
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

	CommandBar components.CommandBar
}

func NewHomeModel() *HomeModel {
	ta := textarea.New()
	ta.Placeholder = "Add notes here..."
	return &HomeModel{
		Problems:   []UserProblemWithRelations{},
		NotesInput: ta,
		CommandBar: components.NewCommandBar([]components.Command{
			{Key: "↑/↓", Desc: "navigate"},
			{Key: "c", Desc: "complete"},
			{Key: "o", Desc: "open"},
			{Key: "?", Desc: "help"},
			{Key: "q", Desc: "quit"},
		}),
	}
}
