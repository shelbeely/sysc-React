package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Alignment constants
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

// RenderOptions holds all configuration for rendering text
type RenderOptions struct {
	Font         *BitFont
	Text         string
	Alignment    int
	Color        string
	Scale        float64
	Shadow       bool
	CharSpacing  int
	LineSpacing  int
	MaxWidth     int // Canvas width for alignment
}

// RenderBitText renders text using a bitmap font with all styling options
func RenderBitText(opts RenderOptions) []string {
	if opts.Font == nil || opts.Text == "" {
		return []string{}
	}

	// Default scale
	if opts.Scale <= 0 {
		opts.Scale = 1.0
	}

	// Step 1: Render base text using font
	lines := opts.Font.RenderText(opts.Text)
	if len(lines) == 0 {
		return lines
	}

	// Step 2: Apply character spacing
	if opts.CharSpacing > 0 {
		lines = applyCharacterSpacing(lines, opts.CharSpacing)
	}

	// Step 3: Apply line spacing
	if opts.LineSpacing > 0 {
		lines = applyLineSpacing(lines, opts.LineSpacing, opts.Font.GetHeight())
	}

	// Step 4: Apply scale
	if opts.Scale != 1.0 {
		lines = applyScale(lines, opts.Scale)
	}

	// Step 5: Apply shadow
	if opts.Shadow {
		lines = applyShadow(lines)
	}

	// Step 6: Apply alignment
	if opts.MaxWidth > 0 {
		lines = applyAlignment(lines, opts.MaxWidth, opts.Alignment)
	}

	// Step 7: Apply color
	if opts.Color != "" {
		lines = applyColor(lines, opts.Color)
	}

	return lines
}

// applyCharacterSpacing adds extra spaces between characters
func applyCharacterSpacing(lines []string, spacing int) []string {
	if spacing <= 0 {
		return lines
	}

	spacer := strings.Repeat(" ", spacing)
	result := make([]string, len(lines))

	for i, line := range lines {
		// Add spacing after each non-space character
		var newLine strings.Builder
		runes := []rune(line)
		for j, r := range runes {
			newLine.WriteRune(r)
			// Add spacer after each character except the last
			if j < len(runes)-1 && r != ' ' {
				newLine.WriteString(spacer)
			}
		}
		result[i] = newLine.String()
	}

	return result
}

// applyLineSpacing adds extra blank lines between text rows
func applyLineSpacing(lines []string, spacing, charHeight int) []string {
	if spacing <= 0 || charHeight <= 0 {
		return lines
	}

	var result []string
	blankLine := ""

	// Calculate number of text blocks (each is charHeight lines)
	numBlocks := (len(lines) + charHeight - 1) / charHeight

	for blockIdx := 0; blockIdx < numBlocks; blockIdx++ {
		start := blockIdx * charHeight
		end := start + charHeight
		if end > len(lines) {
			end = len(lines)
		}

		// Add the block of lines
		result = append(result, lines[start:end]...)

		// Add spacing lines (except after last block)
		if blockIdx < numBlocks-1 {
			// Get line width from first line of block
			if start < len(lines) && len(lines[start]) > 0 {
				blankLine = strings.Repeat(" ", len([]rune(lines[start])))
			}
			for i := 0; i < spacing; i++ {
				result = append(result, blankLine)
			}
		}
	}

	return result
}

// applyScale scales the text by the given factor
func applyScale(lines []string, scale float64) []string {
	if scale == 1.0 {
		return lines
	}

	// For simplicity, we support integer scales for now
	// TODO: Implement fractional scaling (0.5x) using character selection
	intScale := int(scale)
	if intScale < 1 {
		intScale = 1
	}

	var result []string

	// Scale vertically (repeat each line)
	for _, line := range lines {
		// Scale horizontally (repeat each character)
		var scaledLine strings.Builder
		for _, r := range line {
			for i := 0; i < intScale; i++ {
				scaledLine.WriteRune(r)
			}
		}
		scaledLineStr := scaledLine.String()

		// Repeat line vertically
		for i := 0; i < intScale; i++ {
			result = append(result, scaledLineStr)
		}
	}

	return result
}

// applyShadow adds a drop shadow effect
func applyShadow(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		width := len([]rune(line))
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Create shadow offset (1 right, 1 down)
	shadowChar := '░' // Light shade for shadow
	result := make([]string, len(lines)+1)

	// First line has no shadow (nothing above it)
	result[0] = lines[0]

	// Subsequent lines
	for i := 1; i < len(lines); i++ {
		runes := []rune(lines[i])
		prevRunes := []rune(lines[i-1])

		var newLine strings.Builder

		for j := 0; j < maxWidth; j++ {
			var currentChar rune = ' '
			if j < len(runes) {
				currentChar = runes[j]
			}

			// Check if shadow should be drawn here
			// Shadow appears if: current position is space AND previous line had character to the left
			if currentChar == ' ' && j > 0 && j-1 < len(prevRunes) && prevRunes[j-1] != ' ' {
				newLine.WriteRune(shadowChar)
			} else {
				newLine.WriteRune(currentChar)
			}
		}

		result[i] = newLine.String()
	}

	// Add shadow line at bottom
	var shadowLine strings.Builder
	lastRunes := []rune(lines[len(lines)-1])
	for j := 0; j < maxWidth; j++ {
		if j > 0 && j-1 < len(lastRunes) && lastRunes[j-1] != ' ' {
			shadowLine.WriteRune(shadowChar)
		} else {
			shadowLine.WriteRune(' ')
		}
	}
	result[len(lines)] = shadowLine.String()

	return result
}

// applyAlignment aligns text within a given width
func applyAlignment(lines []string, maxWidth, alignment int) []string {
	result := make([]string, len(lines))

	for i, line := range lines {
		lineWidth := len([]rune(line))

		if lineWidth >= maxWidth {
			// Line is too wide, don't modify
			result[i] = line
			continue
		}

		padding := maxWidth - lineWidth

		switch alignment {
		case AlignLeft:
			// No change needed
			result[i] = line + strings.Repeat(" ", padding)

		case AlignCenter:
			leftPad := padding / 2
			rightPad := padding - leftPad
			result[i] = strings.Repeat(" ", leftPad) + line + strings.Repeat(" ", rightPad)

		case AlignRight:
			result[i] = strings.Repeat(" ", padding) + line

		default:
			result[i] = line
		}
	}

	return result
}

// applyColor applies a color to all non-space characters using lipgloss
func applyColor(lines []string, hexColor string) []string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(hexColor))

	result := make([]string, len(lines))
	for i, line := range lines {
		// Apply color to entire line for now
		// In a more sophisticated implementation, we could color character-by-character
		result[i] = style.Render(line)
	}

	return result
}

// GetRenderedDimensions calculates the final dimensions of rendered text
func GetRenderedDimensions(opts RenderOptions) (width, height int) {
	lines := RenderBitText(opts)
	if len(lines) == 0 {
		return 0, 0
	}

	height = len(lines)
	for _, line := range lines {
		w := len([]rune(line))
		if w > width {
			width = w
		}
	}

	return width, height
}
