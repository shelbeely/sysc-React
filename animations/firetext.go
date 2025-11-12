package animations

import (
	"fmt"
	"math/rand"
	"strings"
)

// FireTextEffect implements fire animation with ASCII art displayed as negative space
// Fire burns around the text, creating text shape with empty areas
type FireTextEffect struct {
	width   int      // Terminal width
	height  int      // Terminal height
	buffer  []int    // Heat values (0-65), size = width * height
	palette []string // Hex color codes from theme
	chars   []rune   // Fire characters for density (8-level gradient)

	// Text masking
	text         string
	textMask     [][]bool // [y][x] = true if character exists at this position
	centerX      int
	centerY      int
	artWidth     int
	artHeight    int
}

// NewFireTextEffect creates a new fire-text effect with given dimensions, palette, and ASCII art
func NewFireTextEffect(width, height int, palette []string, text string) *FireTextEffect {
	f := &FireTextEffect{
		width:   width,
		height:  height,
		palette: palette,
		text:    text,
		// Enhanced 8-character gradient for smoother fire rendering
		chars: []rune{' ', '░', '░', '▒', '▒', '▓', '▓', '█'},
	}
	f.parseText()
	f.init()
	return f
}

// parseText extracts ASCII art character positions and creates mask
func (f *FireTextEffect) parseText() {
	lines := strings.Split(f.text, "\n")
	f.artHeight = len(lines)

	// Find max line width
	f.artWidth = 0
	for _, line := range lines {
		if len([]rune(line)) > f.artWidth {
			f.artWidth = len([]rune(line))
		}
	}

	// Center the art
	f.centerX = (f.width - f.artWidth) / 2
	f.centerY = (f.height - f.artHeight) / 2

	// Initialize mask
	f.textMask = make([][]bool, f.height)
	for i := range f.textMask {
		f.textMask[i] = make([]bool, f.width)
	}

	// Mark character positions in mask
	for lineIdx, line := range lines {
		lineRunes := []rune(line)
		for charIdx, char := range lineRunes {
			if char != ' ' && char != '\n' {
				x := f.centerX + charIdx
				y := f.centerY + lineIdx

				// Only mark if within bounds
				if x >= 0 && x < f.width && y >= 0 && y < f.height {
					f.textMask[y][x] = true
				}
			}
		}
	}
}

// Initialize fire buffer with bottom row as heat source
func (f *FireTextEffect) init() {
	f.buffer = make([]int, f.width*f.height)

	// Set bottom row to maximum heat (fire source)
	// Don't place fire where text is
	for i := 0; i < f.width; i++ {
		if !f.textMask[f.height-1][i] {
			f.buffer[(f.height-1)*f.width+i] = 65
		}
	}
}

// UpdatePalette changes the fire color palette (for theme switching)
func (f *FireTextEffect) UpdatePalette(palette []string) {
	f.palette = palette
}

// Resize reinitializes the fire effect with new dimensions
func (f *FireTextEffect) Resize(width, height int) {
	f.width = width
	f.height = height
	f.parseText() // Re-parse to recenter text
	f.init()
}

// spreadFire propagates heat upward with random decay, respecting text mask
func (f *FireTextEffect) spreadFire(from int) {
	fromY := from / f.width
	fromX := from % f.width

	// Don't spread fire from masked positions
	if f.textMask[fromY][fromX] {
		f.buffer[from] = 0 // Ensure masked areas stay cold
		return
	}

	// Random horizontal offset (0-3) for flickering effect
	offset := rand.Intn(4)
	to := from - f.width - offset + 1

	// Bounds check
	if to < 0 || to >= len(f.buffer) {
		return
	}

	toY := to / f.width
	toX := to % f.width

	// Don't spread fire TO masked positions (this creates the negative space)
	if toY >= 0 && toY < f.height && toX >= 0 && toX < f.width {
		if f.textMask[toY][toX] {
			return // Fire cannot enter text areas
		}
	}

	// Random decay (0-3) for natural fade
	decay := rand.Intn(4)

	newHeat := f.buffer[from] - decay
	if newHeat < 0 {
		newHeat = 0
	}

	f.buffer[to] = newHeat
}

// Update advances the fire simulation by one frame
func (f *FireTextEffect) Update() {
	// Process all pixels from bottom to top
	// (Fire spreads upward, must process bottom row first)
	for y := f.height - 1; y > 0; y-- {
		for x := 0; x < f.width; x++ {
			index := y*f.width + x
			f.spreadFire(index)
		}
	}

	// Ensure masked areas stay cold (no fire)
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			if f.textMask[y][x] {
				f.buffer[y*f.width+x] = 0
			}
		}
	}
}

// Render converts fire to colored block output with batched raw ANSI codes
// Text areas are rendered as empty space (negative space effect)
func (f *FireTextEffect) Render() string {
	var output strings.Builder

	// Always render full viewport height to anchor fire at bottom
	for y := 0; y < f.height; y++ {
		var currentColor string
		var batchChars strings.Builder

		for x := 0; x < f.width; x++ {
			heat := f.buffer[y*f.width+x]

			// Text mask areas are always empty (negative space)
			if f.textMask[y][x] {
				// Flush any pending batch
				if batchChars.Len() > 0 {
					r, g, b := hexToRGB(currentColor)
					fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
					batchChars.Reset()
				}
				output.WriteString(" ")
				currentColor = ""
				continue
			}

			// Skip very low heat (natural fade to background)
			if heat < 5 {
				// Flush any pending batch
				if batchChars.Len() > 0 {
					r, g, b := hexToRGB(currentColor)
					fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
					batchChars.Reset()
				}
				output.WriteString(" ")
				currentColor = ""
				continue
			}

			// Map heat to character (0-65 → 8 chars)
			charIndex := (heat * (len(f.chars) - 1)) / 65
			if charIndex >= len(f.chars) {
				charIndex = len(f.chars) - 1
			}
			char := f.chars[charIndex]

			// Map heat to color from palette
			colorIndex := (heat * (len(f.palette) - 1)) / 65
			if colorIndex >= len(f.palette) {
				colorIndex = len(f.palette) - 1
			}
			colorHex := f.palette[colorIndex]

			// If color changed, flush previous batch and start new one
			if colorHex != currentColor {
				if batchChars.Len() > 0 {
					r, g, b := hexToRGB(currentColor)
					fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
					batchChars.Reset()
				}
				currentColor = colorHex
			}

			// Add character to batch
			batchChars.WriteRune(char)
		}

		// Flush any remaining batch at end of line
		if batchChars.Len() > 0 {
			r, g, b := hexToRGB(currentColor)
			fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
		}

		output.WriteString("\n")
	}

	// Remove trailing newline
	result := output.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result
}
