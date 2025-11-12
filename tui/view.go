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

	// If in BIT editor mode, render BIT editor view
	if m.bitEditorMode {
		return m.renderBitEditorView()
	}

	// If in editor mode, render editor view
	if m.editorMode {
		return m.renderEditorView()
	}

	var sections []string

	// Canvas area (viewport for animations)
	sections = append(sections, m.renderCanvas())

	// Selector area
	sections = append(sections, m.renderSelectors())

	// Guidance box (explains current selection)
	sections = append(sections, m.renderGuidance())

	// Help text (no j/k hints)
	sections = append(sections, m.renderHelp())

	// Join all sections and apply background
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
}

// renderCanvas renders the animation preview viewport
func (m Model) renderCanvas() string {
	if m.canvasHeight <= 0 {
		m.canvasHeight = 20 // fallback
	}

	canvasWidth := m.width - 8 // Account for borders and padding

	var content string
	if m.animationRunning && m.currentAnim != nil {
		// Render actual animation frame
		content = m.currentAnim.Render()
	} else {
		// Show welcome/instructions
		content = m.renderWelcome()
	}

	// Pad content to fill canvas
	lines := strings.Split(content, "\n")

	// Truncate or pad lines to match canvas height
	for len(lines) < m.canvasHeight {
		lines = append(lines, "")
	}
	if len(lines) > m.canvasHeight {
		lines = lines[:m.canvasHeight]
	}

	// Pad each line to canvas width (center if possible)
	for i, line := range lines {
		plainLen := len(stripANSI(line))
		if plainLen < canvasWidth {
			// Center the line
			padding := (canvasWidth - plainLen) / 2
			lines[i] = strings.Repeat(" ", padding) + line + strings.Repeat(" ", canvasWidth-plainLen-padding)
		} else if plainLen > canvasWidth {
			lines[i] = line[:canvasWidth]
		}
	}

	canvasContent := strings.Join(lines, "\n")
	return m.styles.Canvas.Width(canvasWidth).Height(m.canvasHeight).Render(canvasContent)
}

// renderWelcome renders the welcome screen
func (m Model) renderWelcome() string {
	welcome := `▄▀▀▀▀ █   █ ▄▀▀▀▀ ▄▀▀▀▀       ▄▀▀▀▀ ▄▀▀▀▄    ▄▀    ▄▀
 ▀▀▀▄ ▀▀▀▀█  ▀▀▀▄ █     ▀▀▀▀▀ █ ▀▀█ █   █  ▄▀    ▄▀
▀▀▀▀  ▀▀▀▀▀ ▀▀▀▀   ▀▀▀▀        ▀▀▀   ▀▀▀  ▀     ▀
             /// SEE YOU SPACE COWBOY//

Terminal Animation Library - Interactive TUI

Select animation settings below
Press ENTER to preview animation in viewport
Press ESC to stop preview`
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

	// Check if this is the File selector and current animation doesn't need a file
	isFileSelector := (index == 2)
	animName := m.animations[m.selectedAnimation]
	needsFile := animName == "matrix-art" || animName == "rain-art" || animName == "pour" ||
		animName == "print" || animName == "beam-text" || animName == "ring-text" || animName == "blackhole-text"

	// Disable file selector for non-text animations
	if isFileSelector && !needsFile {
		style = m.styles.Selector.Faint(true)
		value = "(disabled)"
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

	labelStr := m.styles.SelectorLabel.Render(label)
	valueStr := m.styles.SelectorValue.Render(value)

	// Position indicator with smaller styling
	positionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4C566A")).
		Faint(true)
	positionStr := positionStyle.Render(position)

	content := fmt.Sprintf("%s\n%s\n%s", labelStr, valueStr, positionStr)

	return style.Render(content)
}

// renderHelp renders the help text
func (m Model) renderHelp() string {
	var helpText string
	if m.animationRunning {
		helpText = "ESC Stop animation • ↑/↓ Navigate options • ←/→ Change selector"
	} else {
		helpText = "↑/↓ Navigate options • ←/→ Change selector • ENTER Start animation • Q Quit"
	}
	return m.styles.Help.Render(helpText)
}

// renderGuidance renders guidance/explainer box for current selection
func (m Model) renderGuidance() string {
	animName := m.animations[m.selectedAnimation]
	fileName := m.files[m.selectedFile]

	var guidance string

	// Explain selected animation
	switch animName {
	case "fire":
		guidance = "🔥 DOOM PSX-style fire effect with upward propagation and random flickering"
	case "matrix":
		guidance = "💚 Digital rain with falling character streaks (no text required)"
	case "matrix-art":
		guidance = "💚 Matrix rain effect that reveals your text file content"
	case "rain":
		guidance = "🌧  ASCII character rain effect (no text required)"
	case "rain-art":
		guidance = "🌧  Rain effect that reveals your text file content"
	case "fireworks":
		guidance = "🎆 Physics-based particle fireworks"
	case "pour":
		guidance = "🌊 Text pours down like liquid, character by character"
	case "print":
		guidance = "⌨  Typewriter effect - text appears with typing animation"
	case "beams":
		guidance = "✨ Colored light beams sweep across the screen"
	case "beam-text":
		guidance = "✨ Light beams reveal your text content dramatically"
	case "ring-text":
		guidance = "⭕ Text orbits in 3D rings with perspective effects"
	case "blackhole-text":
		guidance = "🌀 Text gets pulled into a gravitational vortex"
	case "aquarium":
		guidance = "🐠 Animated aquarium with swimming fish and bubbles"
	default:
		guidance = fmt.Sprintf("Selected: %s", animName)
	}

	// Add file info if relevant
	if fileName == "BIT Text Editor" {
		guidance += "\n\n📝 BIT Text Editor: Create ASCII art with 130 fonts"
	} else if fileName == "Custom text" {
		guidance += "\n\n✏  Custom Text: Write or paste your own ASCII art"
	} else if fileName != "(disabled)" {
		guidance += fmt.Sprintf("\n\n📄 Using: %s", fileName)
	}

	return m.styles.GuidanceBox.Render(guidance)
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

	// Apply full background to entire editor view
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
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

	// Apply full background to entire prompt view
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
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

	// Apply full background to entire save prompt view
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
}
