package controllers

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/data/db/dbgen"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
	"github.com/michael-duren/grind-75-cli/internal/utils"
)

func Home(m *models.AppModel, msg tea.Msg) (*models.AppModel, tea.Cmd) {
	slog.Debug("Home Controller Received Message", "msg", msg)

	// Refresh data if needed (naive approach: refresh on every enter for now, but check if we need to re-init table)
	// Ideally we should have a separate message for data refreshing
	if m.Services != nil {
		problems, err := GetUserProblemsWithRelations(m, context.Background())
		if err != nil {
			m.Error = "Failed to get user problems with relations"
			slog.Error("Failed to get user problems with relations", "error", err)
			return m, nil
		}
		m.Home.Problems = problems
	}

	// Initialize table if not done or if we want to refresh rows
	if !m.Home.TableInitialized {
		columns := []table.Column{
			{Title: "Status", Width: 10},
			{Title: "Problem", Width: 30},
			{Title: "Difficulty", Width: 10},
			{Title: "Topic", Width: 15},
			{Title: "Time", Width: 10},
			{Title: "Attempts", Width: 8},
		}

		rows := make([]table.Row, len(m.Home.Problems))
		for i, p := range m.Home.Problems {
			topic := ""
			if len(p.Topics) > 0 {
				topic = p.Topics[0].Name
			}

			rows[i] = table.Row{
				p.Status,
				p.Title,
				p.DifficultyName,
				topic,
				fmt.Sprintf("%d mins", p.Duration),
				fmt.Sprintf("%d", p.Attempts),
			}
		}

		tableHeight := max(m.Height-10, 5)

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(tableHeight),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.ColorTextSub).
			BorderBottom(true).
			Bold(true).
			Foreground(theme.ColorBrand)

		s.Selected = s.Selected.
			Foreground(theme.ColorTextMain).
			Background(theme.ColorBrand).
			Bold(false)

		t.SetStyles(s)

		m.Home.Table = t
		m.Home.TableInitialized = true
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// In Home view, 'q' also quits
			return m, tea.Quit
		case " ":
			// Toggle completion status
			idx := m.Home.Table.Cursor()
			if idx >= 0 && idx < len(m.Home.Problems) {
				p := m.Home.Problems[idx]
				newStatus := "completed"
				if p.Status == "completed" {
					newStatus = "uncompleted" // or "pending" depending on logic, effectively not completed
				}

				// Update DB
				err := m.Services.Queries().UpsertUserProgress(context.Background(), dbgen.UpsertUserProgressParams{
					ProblemID:       p.ProblemID,
					Status:          newStatus,
					LastAttemptedAt: sql.NullTime{Time: time.Now(), Valid: true},
					Attempts:        p.Attempts + 1, // Maybe don't increment attempts on unchecked? debating
				})

				if err != nil {
					slog.Error("Failed to update status", "error", err)
				} else {
					// Optimistic update of UI
					m.Home.Problems[idx].Status = newStatus
					// Re-render rows needed.
					// For now, simpler to mark TableInitialized false to force re-render on next update loop
					// or manually update row. manual update is better.
					rows := m.Home.Table.Rows()
					statusIcon := " "
					if newStatus == "completed" {
						statusIcon = "✅"
					}
					rows[idx][0] = statusIcon
					m.Home.Table.SetRows(rows)
				}
			}
		case "o":
			// Open in browser
			idx := m.Home.Table.Cursor()
			if idx >= 0 && idx < len(m.Home.Problems) {
				url := m.Home.Problems[idx].URL
				return m, openBrowser(url)
			}
		case "right", "l":
			m.Home.SelectedCol = (m.Home.SelectedCol + 1) % 6 // 6 columns
		case "left", "h":
			m.Home.SelectedCol--
			if m.Home.SelectedCol < 0 {
				m.Home.SelectedCol = 5
			}
		}
	case tea.WindowSizeMsg:
		h := max(msg.Height-10, 5)
		m.Home.Table.SetHeight(h)
	}

	m.Home.Table, cmd = m.Home.Table.Update(msg)
	return m, cmd
}

func openBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			slog.Error("Failed to open browser", "error", err)
		}
		return nil
	}
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

func saveNotes(m *models.AppModel, problemID int64, notes string) {
	ctx := context.Background()
	db := m.Services.DB()

	// Find review ID from model
	var reviewID int64
	for _, p := range m.Home.Problems {
		if p.ProblemID == problemID {
			if len(p.Reviews) > 0 {
				reviewID = p.Reviews[0].ID
			}
			break
		}
	}

	if reviewID != 0 {
		// Update existing review
		query := `UPDATE reviews SET notes = ?, review_date = ? WHERE id = ?`
		_, err := db.ExecContext(ctx, query, notes, time.Now(), reviewID)
		if err != nil {
			slog.Error("Failed to update review notes", "error", err)
		}
	} else {
		// Create new review with notes
		query := `INSERT INTO reviews (problem_id, review_date, notes) VALUES (?, ?, ?)`
		_, err := db.ExecContext(ctx, query, problemID, time.Now(), notes)
		if err != nil {
			slog.Error("Failed to create review with notes", "error", err)
		}
	}
}
