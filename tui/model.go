package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI state
type Model struct {
	width  int
	height int

	// Canvas area for animation preview
	canvasHeight int

	// Available options
	animations []string
	themes     []string
	files      []string
	durations  []string

	// Current selections
	selectedAnimation int
	selectedTheme     int
	selectedFile      int
	selectedDuration  int

	// Which selector is focused (0=animation, 1=theme, 2=file, 3=duration)
	focusedSelector int

	// Styles
	styles Styles
}

// Styles holds lipgloss styles for the TUI
type Styles struct {
	Canvas         lipgloss.Style
	Selector       lipgloss.Style
	SelectorFocused lipgloss.Style
	SelectorLabel  lipgloss.Style
	SelectorValue  lipgloss.Style
	Help           lipgloss.Style
}

// NewModel creates a new TUI model with default values
func NewModel() Model {
	// Discover .txt files in assets folder
	files := discoverAssetFiles()
	if len(files) == 0 {
		files = []string{"SYSC.txt"} // fallback
	}

	return Model{
		animations: []string{
			"fire",
			"matrix",
			"matrix-art",
			"rain",
			"rain-art",
			"fireworks",
			"pour",
			"print",
			"beams",
			"beam-text",
			"ring-text",
			"blackhole-text",
			"aquarium",
		},
		themes: []string{
			"dracula",
			"gruvbox",
			"nord",
			"tokyo-night",
			"catppuccin",
			"material",
			"solarized",
			"monochrome",
			"transishardjob",
			"rama",
			"eldritch",
			"dark",
			"default",
		},
		files:             files,
		durations: []string{
			"5s",
			"10s",
			"30s",
			"60s",
			"infinite",
		},
		selectedAnimation: 0,
		selectedTheme:     0,
		selectedFile:      0,
		selectedDuration:  4, // infinite by default
		focusedSelector:   0,
		styles:            NewStyles(),
	}
}

// NewStyles creates the dark theme styles
func NewStyles() Styles {
	return Styles{
		Canvas: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3B4252")).
			Padding(1, 2).
			MarginBottom(1).
			Background(lipgloss.Color("#1E1E2E")),

		Selector: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#3B4252")).
			Padding(0, 2).
			MarginRight(1).
			Background(lipgloss.Color("#1E1E2E")),

		SelectorFocused: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#88C0D0")).
			Padding(0, 2).
			MarginRight(1).
			Background(lipgloss.Color("#2E3440")),

		SelectorLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#88C0D0")).
			Bold(true),

		SelectorValue: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ECEFF4")),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4C566A")).
			MarginTop(1),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// TickMsg is sent when animation should update
type TickMsg time.Time
