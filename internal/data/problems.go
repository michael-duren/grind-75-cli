package data

import (
	_ "embed"
	"encoding/json"
)

type Question struct {
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Duration   int      `json:"duration"`
	EPI        *int     `json:"epi,omitempty"`
	Difficulty string   `json:"difficulty"`
	ID         int      `json:"id"`
	Topic      string   `json:"topic"`
	Routines   []string `json:"routines"`
}

//go:embed problems.json
var defaultQuestionsJSON []byte

func LoadDefaultQuestions() ([]Question, error) {
	var questions []Question
	err := json.Unmarshal(defaultQuestionsJSON, &questions)
	if err != nil {
		return nil, err
	}
	return questions, nil
}
