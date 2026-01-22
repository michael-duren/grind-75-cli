package models

import (
	"time"

	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
)

type ProblemState struct {
	// problem info
	ProblemID  int64
	Slug       string
	Title      string
	Url        string
	Duration   int64
	Difficulty dbgen.DifficultyLevel

	// user progress
	status          string
	lastAttemptedAt *time.Time
	attempts        int64

	Topics []dbgen.Topic
}

type HomeModel struct {
	DefaultProblems []dbgen.Problem
}

func NewHomeModel() *HomeModel {
	return &HomeModel{
		DefaultProblems: []dbgen.Problem{},
	}
}
