package theme

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	ColorBrand      = lipgloss.Color("#f5821f")
	ColorTextMain   = lipgloss.Color("#FFFFFF")
	ColorTextSub    = lipgloss.Color("#B7B7B7")
	ColorBgDark     = lipgloss.Color("#333333")
	ColorBgLeetCode = lipgloss.Color("#262626")

	ColorGreen  = lipgloss.Color("#61BB46") // Easy/Success
	ColorYellow = lipgloss.Color("#FDB827") // Medium/Warning
	ColorRed    = lipgloss.Color("#E03A3E") // Hard/Error
)

var (
	// Base
	Base = lipgloss.NewStyle().Foreground(ColorTextMain)

	// App Header
	AppHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBrand).
			Background(ColorBgDark).
			Padding(1, 0, 1, 1)

	// View Titles
	ViewTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTextMain).
			Padding(0, 0, 1, 0)

	// Form Labels
	Label = lipgloss.NewStyle().
		Foreground(ColorTextSub).
		MarginBottom(0)

	// Inputs
	InputFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBrand).
			Padding(0, 1)

	InputBlurred = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorTextSub).
			Padding(0, 1)

	// Buttons
	BtnNormal = lipgloss.NewStyle().
			Foreground(ColorTextSub)

	BtnActive = lipgloss.NewStyle().
			Foreground(ColorBrand).
			Bold(true).
			SetString("> ") // Prefix with >

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true).
			MarginTop(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorRed).
			Padding(1)
)
