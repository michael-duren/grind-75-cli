package services

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/utils"
)

type ProblemStatus string

const (
	ProblemCompleted  ProblemStatus = "Completed"
	NewProblem        ProblemStatus = "New"
	StrugglingProblem ProblemStatus = "Struggling"
)

type ProblemService interface {
	// GetUserProblemsWithRelations(userID int) ([]Problem, error)
	// ToggleProblemCompletion(userID int, problemID int) error
	SetProblemStatus(problemStatus ProblemStatus, p models.UserProblemWithRelations)
	GetUserProblemsWithRelations() ([]models.UserProblemWithRelations, error)
}

type problemService struct {
	m   *models.AppModel
	ctx context.Context
}

func NewProblemService(m *models.AppModel, ctx context.Context) ProblemService {
	return &problemService{
		m:   m,
		ctx: ctx,
	}
}

func (ps *problemService) SetProblemStatus(problemStatus ProblemStatus, p models.UserProblemWithRelations) {
	m, ctx := ps.m, ps.ctx

	err := m.Services.Queries().UpsertUserProgress(ctx, dbgen.UpsertUserProgressParams{
		ProblemID:       p.ProblemID,
		Status:          string(problemStatus),
		LastAttemptedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Attempts:        p.Attempts,
	})

	if err != nil {
		slog.Error("Failed to update status", "error", err)
		return
	}
}

func (ps *problemService) GetUserProblemsWithRelations() ([]models.UserProblemWithRelations, error) {
	m, ctx := ps.m, ps.ctx
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
