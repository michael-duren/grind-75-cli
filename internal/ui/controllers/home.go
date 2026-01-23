package controllers

import (
	"context"
	"log/slog"
	"time"

	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/michael-duren/grind-75-cli/internal/ui/models"
	"github.com/michael-duren/grind-75-cli/internal/ui/services"
	"github.com/michael-duren/grind-75-cli/internal/ui/theme"
)

type HomeController struct {
	model *models.AppModel
}

func NewHomeController(m *models.AppModel) Controller {
	return &HomeController{
		model: m,
	}
}

func (h *HomeController) View() string {
	return ""
}

func (h *HomeController) Update(msg tea.Msg) tea.Cmd {
	m := h.model
	slog.Debug("Home Controller Received Message", "msg", msg)
	// TODO: come back and look more into context considerations
	ctx := context.Background()
	problemService := services.NewProblemService(m, ctx)

	// Refresh data if needed (naive approach: refresh on every enter for now, but check if we need to re-init table)
	// Ideally we should have a separate message for data refreshing
	// TODO: look into this
	if m.Services != nil {
		problems, err := problemService.GetUserProblemsWithRelations()
		if err != nil {
			m.Error = "Failed to get user problems with relations"
			slog.Error("Failed to get user problems with relations", "error", err)
			return nil
		}
		m.Home.Problems = problems
	}

	if !m.Home.TableInitialized {
		intializeTable(m)
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			// In Home view, 'q' also quits
			return tea.Quit
		case "c":
			idx := m.Home.Table.Cursor()
			if idx < 0 || idx > len(m.Home.Problems) {
				m.Error = "index was out of bounds in home table"
				return nil
			}
			p := m.Home.Problems[idx]
			problemService.SetProblemStatus(services.ProblemCompleted, p)

			m.Home.Problems[idx].Status = string(services.ProblemCompleted)
			rows := m.Home.Table.Rows()

			rows[idx][0] = getStatusIcon(services.ProblemCompleted)
			m.Home.Table.SetRows(rows)
			// toggleComplete(m)
		case "o":
			idx := m.Home.Table.Cursor()
			if idx >= 0 && idx < len(m.Home.Problems) {
				url := m.Home.Problems[idx].URL
				return openBrowser(url)
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
		h, v := theme.Base.GetFrameSize()
		m.Width = msg.Width
		m.Height = msg.Height

		// Layout logic for Table Width
		totalWidth := msg.Width
		tableWidth := totalWidth - h

		// Breakpoint 100 for split view
		if totalWidth >= 100 {
			// Split mode: 65% for table
			splitWidth := max(totalWidth*65/100, 65)
			tableWidth = splitWidth - h
		} else {
			// Stacked mode
			tableWidth = totalWidth - h - 4
		}

		m.Home.Table.SetWidth(tableWidth)

		tableHeight := max(msg.Height-v-4, 5)
		m.Home.Table.SetHeight(tableHeight)
	}

	m.Home.Table, cmd = m.Home.Table.Update(msg)
	return cmd
}

func getStatusIcon(status services.ProblemStatus) string {
	switch status {
	case services.ProblemCompleted:
		return "✅"
	case services.StrugglingProblem:
		return "⚠️"
	case services.NewProblem:
		return "🆕"
	default:
		return " "
	}
}

func intializeTable(m *models.AppModel) {
	columns := []table.Column{
		{Title: "Sts", Width: 4},
		{Title: "Problem", Width: 32},
		{Title: "Diff", Width: 8},
		{Title: "Topic", Width: 18},
		{Title: "Time", Width: 12},
		{Title: "Att", Width: 4},
	}

	rows := make([]table.Row, len(m.Home.Problems))
	for i, p := range m.Home.Problems {
		topic := ""
		if len(p.Topics) > 0 {
			topic = p.Topics[0].Name
		}

		rows[i] = table.Row{
			getStatusIcon(services.ProblemStatus(p.Status)),
			p.Title,
			p.DifficultyName,
			topic,
			fmt.Sprintf("%d mins", p.Duration),
			fmt.Sprintf("%d", p.Attempts),
		}
	}

	tableHeight := max(m.Height-20, 5)

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
