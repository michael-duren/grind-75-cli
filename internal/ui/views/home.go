package views

import (
	"log/slog"
	"strings"

	"github.com/michael-duren/grind-75-cli/internal/ui/models"
)

func Home(m *models.AppModel) string {
	sb := strings.Builder{}

	slog.Debug("Rendering Home View", "problems_count", len(m.Home.Problems))
	for _, problem := range m.Home.Problems {
		sb.WriteString("- ")
		sb.WriteString(problem.Title)
		sb.WriteString(" (")
		sb.WriteString(problem.DifficultyName)
		sb.WriteString(")\n")
	}

	if sb.Len() == 0 {
		sb.WriteString("No problems found.\n")
	}

	sb.WriteString("\n(Press s for Settings)")

	return sb.String()
}
