package models

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
)

// type ProblemState struct {
// 	// problem info
// 	ProblemID  int64
// 	Slug       string
// 	Title      string
// 	Url        string
// 	Duration   int64
// 	Difficulty dbgen.DifficultyLevel
//
// 	// user progress
// 	status          string
// 	lastAttemptedAt *time.Time
// 	attempts        int64
//
// 	Topics []dbgen.Topic
// }

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
}

func NewHomeModel() *HomeModel {
	return &HomeModel{
		Problems: []UserProblemWithRelations{},
	}
}
