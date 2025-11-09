package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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

	// Editor mode for custom text creation
	editorMode       bool
	textarea         textarea.Model
	filenameInput    textinput.Model
	showSavePrompt   bool
	showExportPrompt bool
	exportTarget     int    // 0=syscgo, 1=sysc-walls
	saveError        string // Error message from save operation
	savingInProgress bool

	// Styles
	styles Styles
}

// Styles holds lipgloss styles for the TUI
type Styles struct {
	Canvas          lipgloss.Style
	Selector        lipgloss.Style
	SelectorFocused lipgloss.Style
	SelectorLabel   lipgloss.Style
	SelectorValue   lipgloss.Style
	Help            lipgloss.Style
	Background      lipgloss.Style
}

// NewModel creates a new TUI model with default values
func NewModel() Model {
	// Discover .txt files in assets folder
	files := discoverAssetFiles()
	if len(files) == 0 {
		files = []string{"SYSC.txt"} // fallback
	}
	// Prepend "Custom text" option to files list
	files = append([]string{"Custom text"}, files...)

	// Initialize textarea for editor mode
	ta := textarea.New()
	ta.Placeholder = "Enter your ASCII art here..."
	ta.Focus()
	ta.CharLimit = 10000
	ta.SetWidth(100)
	ta.SetHeight(20)
	ta.ShowLineNumbers = true

	// Initialize filename input
	fi := textinput.New()
	fi.Placeholder = "filename.txt"
	fi.CharLimit = 256
	fi.Width = 40

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
		selectedFile:      1, // Default to SYSC.txt (index 1 after prepending "Custom text")
		selectedDuration:  4, // infinite by default
		focusedSelector:   0,
		editorMode:        false,
		textarea:          ta,
		filenameInput:     fi,
		showSavePrompt:    false,
		showExportPrompt:  false,
		exportTarget:      0,
		saveError:         "",
		savingInProgress:  false,
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
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3B4252")).
			Padding(0, 1).
			MarginRight(1).
			Width(16).
			Background(lipgloss.Color("#1E1E2E")),

		SelectorFocused: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#88C0D0")).
			Padding(0, 1).
			MarginRight(1).
			Width(16).
			Background(lipgloss.Color("#2E3440")),

		SelectorLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#88C0D0")).
			Bold(true).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#3B4252")).
			PaddingBottom(0).
			MarginBottom(0),

		SelectorValue: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ECEFF4")).
			Bold(false).
			PaddingTop(0).
			MarginTop(0),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4C566A")).
			Background(lipgloss.Color("#1E1E2E")).
			Padding(1, 2),

		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#1E1E2E")),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// TickMsg is sent when animation should update
type TickMsg time.Time
