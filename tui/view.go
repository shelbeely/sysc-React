package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// If in editor mode, render editor view
	if m.editorMode {
		return m.renderEditorView()
	}

	var sections []string

	// Canvas area
	sections = append(sections, m.renderCanvas())

	// Selector area
	sections = append(sections, m.renderSelectors())

	// Help text
	sections = append(sections, m.renderHelp())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderCanvas renders the animation preview canvas
func (m Model) renderCanvas() string {
	if m.canvasHeight <= 0 {
		m.canvasHeight = 20 // fallback
	}

	canvasWidth := m.width - 8 // Account for borders and padding

	// Always show welcome screen (no preview mode since we launch immediately)
	content := m.renderWelcome()

	// Pad content to fill canvas height
	lines := strings.Split(content, "\n")
	for len(lines) < m.canvasHeight {
		lines = append(lines, "")
	}
	if len(lines) > m.canvasHeight {
		lines = lines[:m.canvasHeight]
	}

	// Pad each line to canvas width
	for i, line := range lines {
		if len(line) < canvasWidth {
			lines[i] = line + strings.Repeat(" ", canvasWidth-len(line))
		} else if len(line) > canvasWidth {
			lines[i] = line[:canvasWidth]
		}
	}

	canvasContent := strings.Join(lines, "\n")
	return m.styles.Canvas.Width(canvasWidth).Height(m.canvasHeight).Render(canvasContent)
}

// renderWelcome renders the welcome screen
func (m Model) renderWelcome() string {
	welcome := `
 ███████╗██╗   ██╗███████╗ ██████╗
 ██╔════╝╚██╗ ██╔╝██╔════╝██╔════╝
 ███████╗ ╚████╔╝ ███████╗██║
 ╚════██║  ╚██╔╝  ╚════██║██║
 ███████║   ██║   ███████║╚██████╗
 ╚══════╝   ╚═╝   ╚══════╝ ╚═════╝

 Terminal Animation Library - TUI

 Select animation settings below
 Press ENTER to start animation
`
	return welcome
}

// renderSelectors renders the selector controls
func (m Model) renderSelectors() string {
	selectors := []string{
		m.renderSelector(0, "Animation", m.animations[m.selectedAnimation]),
		m.renderSelector(1, "Theme", m.themes[m.selectedTheme]),
		m.renderSelector(2, "File", m.files[m.selectedFile]),
		m.renderSelector(3, "Duration", m.durations[m.selectedDuration]),
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, selectors...)
}

// renderSelector renders a single selector
func (m Model) renderSelector(index int, label, value string) string {
	style := m.styles.Selector
	if index == m.focusedSelector {
		style = m.styles.SelectorFocused
	}

	// Calculate position indicator
	var position string
	switch index {
	case 0: // Animation
		position = fmt.Sprintf("(%d/%d)", m.selectedAnimation+1, len(m.animations))
	case 1: // Theme
		position = fmt.Sprintf("(%d/%d)", m.selectedTheme+1, len(m.themes))
	case 2: // File
		position = fmt.Sprintf("(%d/%d)", m.selectedFile+1, len(m.files))
	case 3: // Duration
		position = fmt.Sprintf("(%d/%d)", m.selectedDuration+1, len(m.durations))
	}

	labelStr := m.styles.SelectorLabel.Render(label + " ▼")
	valueStr := m.styles.SelectorValue.Render(value)
	positionStr := m.styles.Help.Render(position)

	content := fmt.Sprintf("%s\n%s %s", labelStr, valueStr, positionStr)

	// Add navigation hint for focused selector
	if index == m.focusedSelector {
		content += "\n" + m.styles.Help.Render("↑↓")
	}

	return style.Render(content)
}

// renderHelp renders the help text
func (m Model) renderHelp() string {
	helpText := "↑/↓ or j/k Navigate • ←/→ or h/l Change selector • Enter Start animation • Esc Quit"
	return m.styles.Help.Render(helpText)
}

// renderEditorView renders the ASCII text editor
func (m Model) renderEditorView() string {
	if m.showExportPrompt {
		return m.renderExportPrompt()
	}

	if m.showSavePrompt {
		return m.renderSavePrompt()
	}

	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("ASCII Text Editor")
	sections = append(sections, title)

	// Textarea
	textareaStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(1, 2).
		Width(m.width - 6).
		Background(lipgloss.Color("#1E1E2E"))

	sections = append(sections, textareaStyle.Render(m.textarea.View()))

	// Help text
	helpText := "Type your ASCII art • Ctrl+S Save/Export • Esc Cancel"
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4C566A")).
		Padding(1, 0)
	sections = append(sections, helpStyle.Render(helpText))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderExportPrompt renders the export target selection dialog
func (m Model) renderExportPrompt() string {
	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("Save/Export ASCII Art")
	sections = append(sections, title)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Padding(1, 0).
		Render("Select where to save:")
	sections = append(sections, instructions)

	// Export options
	exportOptions := []string{
		"syscgo - Save to assets/ folder for animations",
		"sysc-walls (WIP) - Save as wallpaper (coming soon)",
	}

	optionStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(lipgloss.Color("#ECEFF4"))

	selectedStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Background(lipgloss.Color("#2E3440"))

	var optionsRendered []string
	for i, option := range exportOptions {
		if i == m.exportTarget {
			optionsRendered = append(optionsRendered, selectedStyle.Render("▸ "+option))
		} else {
			optionsRendered = append(optionsRendered, optionStyle.Render("  "+option))
		}
	}

	optionsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(1, 2).
		Width(m.width - 6).
		Background(lipgloss.Color("#1E1E2E"))

	sections = append(sections, optionsBox.Render(lipgloss.JoinVertical(lipgloss.Left, optionsRendered...)))

	// Help text
	helpText := "↑/↓ Select • Enter Confirm • Esc Cancel"
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4C566A")).
		Padding(1, 0)
	sections = append(sections, helpStyle.Render(helpText))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderSavePrompt renders the save dialog
func (m Model) renderSavePrompt() string {
	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("Save ASCII Art")
	sections = append(sections, title)

	// Error message if any
	if m.saveError != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BF616A")).
			Bold(true).
			Padding(1, 0)
		sections = append(sections, errorStyle.Render("⚠ "+m.saveError))
	}

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Padding(1, 0).
		Render("Enter filename (will be saved to assets/ folder):")
	sections = append(sections, instructions)

	// Filename input
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(1, 2).
		Width(m.width - 6).
		Background(lipgloss.Color("#2E3440"))
	sections = append(sections, inputStyle.Render(m.filenameInput.View()))

	// Help text
	helpText := "Enter Confirm • Esc Cancel"
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4C566A")).
		Padding(1, 0)
	sections = append(sections, helpStyle.Render(helpText))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
