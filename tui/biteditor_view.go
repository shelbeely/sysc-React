package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderBitEditorView renders the BIT text editor interface
func (m Model) renderBitEditorView() string {
	if m.bitShowFontList {
		return m.renderFontBrowser()
	}

	if m.bitColorPicker {
		return m.renderColorPicker()
	}

	if m.showExportPrompt {
		return m.renderExportPrompt() // Reuse existing export prompt
	}

	if m.showSavePrompt {
		return m.renderBitSavePrompt()
	}

	var sections []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("BIT Text Editor - Banner Text Generator")
	sections = append(sections, title)

	// Preview canvas
	sections = append(sections, m.renderBitPreview())

	// Text input
	sections = append(sections, m.renderBitTextInput())

	// Controls
	sections = append(sections, m.renderBitControls())

	// Help text
	sections = append(sections, m.renderBitHelp())

	// Apply full background
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
}

// renderBitPreview renders the live preview canvas
func (m Model) renderBitPreview() string {
	canvasWidth := m.width - 8
	canvasHeight := m.height - 25 // Leave room for controls
	if canvasHeight < 10 {
		canvasHeight = 10
	}

	var preview string
	if len(m.bitPreviewLines) > 0 {
		// Take first N lines that fit
		displayLines := m.bitPreviewLines
		if len(displayLines) > canvasHeight {
			displayLines = displayLines[:canvasHeight]
		}
		preview = strings.Join(displayLines, "\n")
	} else {
		// Show placeholder
		preview = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4C566A")).
			Render("Preview will appear here... Type text below to see it rendered.")
	}

	return m.styles.Canvas.
		Width(canvasWidth).
		Height(canvasHeight).
		Render(preview)
}

// renderBitTextInput renders the text input field
func (m Model) renderBitTextInput() string {
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(0, 1).
		Width(m.width - 8)

	if m.bitFocusedControl == 0 {
		inputStyle = inputStyle.BorderForeground(lipgloss.Color("#A3BE8C"))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Text: ")

	return inputStyle.Render(label + m.bitTextInput.View())
}

// renderBitControls renders all control panels
func (m Model) renderBitControls() string {
	controlsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3B4252")).
		Padding(1, 2).
		Width(m.width - 8).
		Background(lipgloss.Color("#1E1E2E"))

	var controls []string

	// Row 1: Font, Alignment, Color
	row1 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderFontControl(),
		m.renderAlignmentControl(),
		m.renderColorControl(),
	)
	controls = append(controls, row1)

	// Row 2: Scale, Shadow, Spacing
	row2 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderScaleControl(),
		m.renderShadowControl(),
		m.renderSpacingControl(),
	)
	controls = append(controls, row2)

	return controlsStyle.Render(lipgloss.JoinVertical(lipgloss.Left, controls...))
}

// renderFontControl renders the font selector
func (m Model) renderFontControl() string {
	focused := m.bitFocusedControl == 1
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(20)

	if focused {
		style = style.Background(lipgloss.Color("#2E3440"))
	}

	fontName := "none"
	if m.bitCurrentFont != nil {
		fontName = m.bitCurrentFont.Name
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Font: ")

	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Render(fmt.Sprintf("%s (%d/%d)", fontName, m.bitSelectedFont+1, len(m.bitFonts)))

	return style.Render(label + "\n" + value)
}

// renderAlignmentControl renders alignment buttons
func (m Model) renderAlignmentControl() string {
	focused := m.bitFocusedControl == 2
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(18)

	if focused {
		style = style.Background(lipgloss.Color("#2E3440"))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Align: ")

	buttons := []string{}
	alignments := []string{"[L]", "[C]", "[R]"}
	for i, text := range alignments {
		if i == m.bitAlignment {
			buttons = append(buttons, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A3BE8C")).
				Bold(true).
				Render(text))
		} else {
			buttons = append(buttons, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4C566A")).
				Render(text))
		}
	}

	return style.Render(label + "\n" + strings.Join(buttons, " "))
}

// renderColorControl renders color selector
func (m Model) renderColorControl() string {
	focused := m.bitFocusedControl == 3
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(20)

	if focused {
		style = style.Background(lipgloss.Color("#2E3440"))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Color: ")

	// Show color swatch
	swatch := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.bitColor)).
		Render("███ ")

	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Render(m.bitColor)

	return style.Render(label + "\n" + swatch + value)
}

// renderScaleControl renders scale selector
func (m Model) renderScaleControl() string {
	focused := m.bitFocusedControl == 4
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(18)

	if focused {
		style = style.Background(lipgloss.Color("#2E3440"))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Scale: ")

	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Render(fmt.Sprintf("%.1fx", m.bitScale))

	return style.Render(label + "\n" + value)
}

// renderShadowControl renders shadow toggle
func (m Model) renderShadowControl() string {
	focused := m.bitFocusedControl == 5
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(20)

	if focused {
		style = style.Background(lipgloss.Color("#2E3440"))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Shadow: ")

	status := "Off"
	if m.bitShadow {
		status = fmt.Sprintf("On (%d,%d)", m.bitShadowOffsetX, m.bitShadowOffsetY)
	}

	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Render(status)

	return style.Render(label + "\n" + value)
}

// renderSpacingControl renders spacing controls
func (m Model) renderSpacingControl() string {
	focused := m.bitFocusedControl == 6
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(20)

	if focused {
		style = style.Background(lipgloss.Color("#2E3440"))
	}

	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88C0D0")).
		Render("Spacing: ")

	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Render(fmt.Sprintf("C:%d W:%d L:%d", m.bitCharSpacing, m.bitWordSpacing, m.bitLineSpacing))

	return style.Render(label + "\n" + value)
}

// renderBitHelp renders help text for BIT editor
func (m Model) renderBitHelp() string {
	helpText := "Tab/Shift+Tab Controls • ←/→ Adjust • Enter Select • F Font List • C Color • Ctrl+S Save • Esc Back"
	return m.styles.Help.Render(helpText)
}

// renderFontBrowser renders the font selection browser
func (m Model) renderFontBrowser() string {
	var sections []string

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("Select Font")
	sections = append(sections, title)

	// Font list
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(1, 2).
		Width(m.width - 8).
		Height(m.height - 10).
		Background(lipgloss.Color("#1E1E2E"))

	var fontItems []string
	startIdx := 0
	endIdx := len(m.bitFonts)

	// Show window of fonts around selection
	windowSize := m.height - 12
	if windowSize < 10 {
		windowSize = 10
	}

	if len(m.bitFonts) > windowSize {
		startIdx = m.bitSelectedFont - windowSize/2
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + windowSize
		if endIdx > len(m.bitFonts) {
			endIdx = len(m.bitFonts)
			startIdx = endIdx - windowSize
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	for i := startIdx; i < endIdx; i++ {
		fontName := m.bitFonts[i]
		if i == m.bitSelectedFont {
			fontItems = append(fontItems, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A3BE8C")).
				Bold(true).
				Render(fmt.Sprintf("▸ %s", fontName)))
		} else {
			fontItems = append(fontItems, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ECEFF4")).
				Render(fmt.Sprintf("  %s", fontName)))
		}
	}

	sections = append(sections, listStyle.Render(strings.Join(fontItems, "\n")))

	helpText := "↑/↓ Navigate • Enter Select • Esc Cancel"
	sections = append(sections, m.styles.Help.Render(helpText))

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
}

// renderColorPicker renders the color picker
func (m Model) renderColorPicker() string {
	var sections []string

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("Select Color")
	sections = append(sections, title)

	// Theme colors
	themeColors := []struct {
		name  string
		color string
	}{
		{"Nord Blue", "#88C0D0"},
		{"Nord Green", "#A3BE8C"},
		{"Nord Purple", "#B48EAD"},
		{"Nord Orange", "#D08770"},
		{"Nord Red", "#BF616A"},
		{"Nord Yellow", "#EBCB8B"},
		{"Dracula Purple", "#BD93F9"},
		{"Dracula Pink", "#FF79C6"},
		{"Dracula Cyan", "#8BE9FD"},
		{"Dracula Green", "#50FA7B"},
		{"White", "#FFFFFF"},
		{"Gray", "#808080"},
	}

	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(1, 2).
		Width(m.width - 8).
		Background(lipgloss.Color("#1E1E2E"))

	var colorItems []string
	for _, c := range themeColors {
		swatch := lipgloss.NewStyle().
			Foreground(lipgloss.Color(c.color)).
			Render("███ ")

		item := swatch + c.name + " " + c.color

		if c.color == m.bitColor {
			colorItems = append(colorItems, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A3BE8C")).
				Bold(true).
				Render("▸ "+item))
		} else {
			colorItems = append(colorItems, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ECEFF4")).
				Render("  "+item))
		}
	}

	sections = append(sections, listStyle.Render(strings.Join(colorItems, "\n")))

	helpText := "↑/↓ Navigate • Enter Select • Esc Cancel"
	sections = append(sections, m.styles.Help.Render(helpText))

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
}

// renderBitSavePrompt renders the save dialog for BIT editor
func (m Model) renderBitSavePrompt() string {
	var sections []string

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#88C0D0")).
		Padding(1, 0).
		Render("Save Banner Text")
	sections = append(sections, title)

	if m.saveError != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BF616A")).
			Bold(true).
			Padding(1, 0)
		sections = append(sections, errorStyle.Render("⚠ "+m.saveError))
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECEFF4")).
		Padding(1, 0).
		Render("Enter filename (will be saved to assets/ folder):")
	sections = append(sections, instructions)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#88C0D0")).
		Padding(1, 2).
		Width(m.width - 6).
		Background(lipgloss.Color("#2E3440"))
	sections = append(sections, inputStyle.Render(m.filenameInput.View()))

	helpText := "Enter Confirm • Esc Cancel"
	sections = append(sections, m.styles.Help.Render(helpText))

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return m.styles.Background.
		Width(m.width).
		Height(m.height).
		Render(content)
}
