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

	var content string
	if m.showPreview {
		content = m.previewContent
	} else {
		// Show welcome screen
		content = m.renderWelcome()
	}

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
 Press ENTER to preview
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

	labelStr := m.styles.SelectorLabel.Render(label)
	valueStr := m.styles.SelectorValue.Render(value)

	content := fmt.Sprintf("%s\n%s", labelStr, valueStr)

	// Add navigation hint for focused selector
	if index == m.focusedSelector {
		content += "\n" + m.styles.Help.Render("↑↓")
	}

	return style.Render(content)
}

// renderHelp renders the help text
func (m Model) renderHelp() string {
	var helpText string
	if m.showPreview {
		helpText = "Enter Launch animation • Esc Go back • Ctrl+C Quit"
	} else {
		helpText = "↑/↓ Navigate • ←/→ Change selector • Enter Preview • Esc Quit"
	}
	return m.styles.Help.Render(helpText)
}
