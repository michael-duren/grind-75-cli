package controllers

import (
	"context"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/utils"
)

func Home(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	slog.Debug("Home Controller Received Message", "msg", msg)

	problems, err := GetUserProblemsWithRelations(m, context.Background())
	if err != nil {
		m.Error = "Failed to get user problems with relations"
		slog.Error("Failed to get user problems with relations", "error", err)
		return m, nil
	}
	m.Home.Problems = problems

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// In Home view, 'q' also quits
			return m, tea.Quit
		}
	}
	return m, nil
}

func GetUserProblemsWithRelations(m *models.AppModel, ctx context.Context) ([]models.UserProblemWithRelations, error) {
	// Get all problems
	queries := m.Services.Queries()
	problems, err := queries.ListUserProblems(ctx)
	if err != nil {
		return nil, err
	}

	// Get all reviews
	reviews, err := queries.GetAllProblemReviews(ctx)
	if err != nil {
		return nil, err
	}

	// Get all topics
	topicRows, err := queries.GetAllProblemTopics(ctx)
	if err != nil {
		return nil, err
	}

	// Group reviews by problem_id
	reviewMap := make(map[int64][]dbgen.Review)
	for _, r := range reviews {
		reviewMap[r.ProblemID] = append(reviewMap[r.ProblemID], dbgen.Review{
			ID:         r.ID,
			ReviewDate: r.ReviewDate,
			Completed:  r.Completed,
			Notes:      r.Notes,
		})
	}

	// Group topics by problem_id
	topicMap := make(map[int64][]dbgen.Topic)
	for _, t := range topicRows {
		topicMap[t.ProblemID] = append(topicMap[t.ProblemID], dbgen.Topic{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	// Combine everything
	result := make([]models.UserProblemWithRelations, len(problems))
	for i, p := range problems {
		result[i] = models.UserProblemWithRelations{
			ProblemID:       p.ProblemID,
			Slug:            p.Slug,
			Title:           p.Title,
			URL:             p.Url,
			Duration:        p.Duration,
			DifficultyID:    utils.CoerceFromNullString(p.DifficultyID),
			DifficultyName:  utils.CoerceFromNullString(p.DifficultyName),
			Status:          utils.CoerceFromNullString(p.Status),
			LastAttemptedAt: utils.CoerceFromNullTime(p.LastAttemptedAt),
			Attempts:        utils.CoerceFromNullInt64(p.Attempts),
			Topics:          topicMap[p.ProblemID],  // Attach topics here
			Reviews:         reviewMap[p.ProblemID], // Attach reviews here
		}

		if result[i].Topics == nil {
			result[i].Topics = []dbgen.Topic{}
		}
		if result[i].Reviews == nil {
			result[i].Reviews = []dbgen.Review{}
		}
	}

	return result, nil
}
