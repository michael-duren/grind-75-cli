package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/michael-duren/grind-75-cli/internal/data"
	dbgen "github.com/michael-duren/grind-75-cli/internal/data/db/gen"
)

// Seed populates the database with the default Grind 75 questions.
func Seed(ctx context.Context, db *sql.DB) error {
	queries := dbgen.New(db)

	questions, err := data.LoadDefaultQuestions()
	if err != nil {
		return fmt.Errorf("failed to load default questions: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	// Note: Difficulty levels are expected to be pre-filled by migrations/schema
	// but strictly speaking we could ensure they exist here if we wanted to be safe.
	// For now we assume they exist.

	for _, q := range questions {
		// Ensure difficulty is valid (simple map check or just insert)
		// We trust the JSON matches 'Easy', 'Medium', 'Hard'

		err := qtx.CreateProblem(ctx, dbgen.CreateProblemParams{
			ID:           int64(q.ID),
			Slug:         q.Slug,
			Title:        q.Title,
			Url:          q.URL,
			Duration:     int64(q.Duration),
			DifficultyID: q.Difficulty,
		})
		if err != nil {
			return fmt.Errorf("failed to create problem %s: %w", q.Slug, err)
		}

		// Handle Topics
		// The JSON has "topic" (string) but some have "routines" (array of strings)
		// We should treat both as topics to link.
		// Actually, let's look at schema: problems has topics.

		// Insert primary topic
		if q.Topic != "" {
			err = qtx.CreateTopic(ctx, dbgen.CreateTopicParams{
				ID:   q.Topic,
				Name: q.Topic,
			})
			if err != nil {
				return fmt.Errorf("failed to create topic %s: %w", q.Topic, err)
			}
			err = qtx.LinkProblemTopic(ctx, dbgen.LinkProblemTopicParams{
				ProblemID: int64(q.ID),
				TopicID:   q.Topic,
			})
			if err != nil {
				return fmt.Errorf("failed to link topic %s: %w", q.Topic, err)
			}
		}

		// Insert routines as topics too?
		for _, routine := range q.Routines {
			err = qtx.CreateTopic(ctx, dbgen.CreateTopicParams{
				ID:   routine,
				Name: routine,
			})
			if err != nil {
				return fmt.Errorf("failed to create routine topic %s: %w", routine, err)
			}
			err = qtx.LinkProblemTopic(ctx, dbgen.LinkProblemTopicParams{
				ProblemID: int64(q.ID),
				TopicID:   routine,
			})
			if err != nil {
				return fmt.Errorf("failed to link routine topic %s: %w", routine, err)
			}
		}

		// Initialize user progress as 'New'
		// We use UpsertUserProgress but that takes params.
		// Wait, user_progress might be empty initially.
		// We can insert with default status 'New'.
		// The queries.sql had UpsertUserProgress. Let's use it or add a dbgeneric init.
		// Actually, we probably don't need to seed user_progress for ALL problems immediately,
		// or maybe we do to simplify querying "New" problems.
		// Let's seed it as 'New'.

		err = qtx.UpsertUserProgress(ctx, dbgen.UpsertUserProgressParams{
			ProblemID:       int64(q.ID),
			Status:          "New",
			LastAttemptedAt: sql.NullTime{}, // Valid: false
			Attempts:        0,
		})
		if err != nil {
			return fmt.Errorf("failed to init user progress for %s: %w", q.Slug, err)
		}
	}

	return tx.Commit()
}
