package animations

import (
	"fmt"
	"math/rand"
	"strings"
)

// FireEffect implements PSX DOOM-style fire algorithm
type FireEffect struct {
	width   int      // Terminal width
	height  int      // Terminal height
	buffer  []int    // Heat values (0-36), size = width * height
	palette []string // Hex color codes from theme
	chars   []rune   // Fire characters for density
}

// NewFireEffect creates a new fire effect with given dimensions and theme palette
func NewFireEffect(width, height int, palette []string) *FireEffect {
	f := &FireEffect{
		width:   width,
		height:  height,
		palette: palette,
		chars:   []rune{' ', '░', '▒', '▓', '█'},
	}
	f.init()
	return f
}

// Initialize fire buffer with bottom row as heat source
func (f *FireEffect) init() {
	f.buffer = make([]int, f.width*f.height)

	// Set bottom row to varied heat (less dense)
	for i := 0; i < f.width; i++ {
		// Random heat 24-36 for less density
		f.buffer[(f.height-1)*f.width+i] = 24 + rand.Intn(13)
	}
}

// UpdatePalette changes the fire color palette (for theme switching)
func (f *FireEffect) UpdatePalette(palette []string) {
	f.palette = palette
}

// Resize reinitializes the fire effect with new dimensions
func (f *FireEffect) Resize(width, height int) {
	f.width = width
	f.height = height
	f.init()
}

// spreadFire propagates heat upward with random decay
func (f *FireEffect) spreadFire(from int) {
	// Random horizontal offset (0-3) for chaos
	offset := rand.Intn(4)
	to := from - f.width - offset + 1

	// Bounds check
	if to < 0 || to >= len(f.buffer) {
		return
	}

	// Calculate target row
	toY := to / f.width
	hardLimit := (f.height * 3) / 10 // Top 30% - absolute no-go zone
	fadeZoneStart := f.height / 2     // Top 50% - gentle fade zone

	// Hard limit - no propagation into top 30%
	if toY < hardLimit {
		return
	}

	// Random decay (0-2) - increased for less density
	decay := rand.Intn(3)

	// Fade zone (between 30% and 50% from top)
	if toY < fadeZoneStart {
		decay += rand.Intn(2) + 1 // Add 1 or 2 extra
	}

	newHeat := f.buffer[from] - decay
	if newHeat < 0 {
		newHeat = 0
	}

	f.buffer[to] = newHeat
}

// Update advances the fire simulation by one frame
func (f *FireEffect) Update() {
	// Process all pixels from bottom to top
	// (Fire spreads upward, must process bottom row first)
	for y := f.height - 1; y > 0; y-- {
		for x := 0; x < f.width; x++ {
			index := y*f.width + x
			f.spreadFire(index)
		}
	}
}

// hexToRGB converts hex color to RGB values
func hexToRGB(hex string) (int, int, int) {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Parse RGB
	var r, g, b int
	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	}
	return r, g, b
}

// Render converts fire to colored block output with batched raw ANSI codes
func (f *FireEffect) Render() string {
	var output strings.Builder

	for y := 0; y < f.height; y++ {
		var currentColor string
		var batchChars strings.Builder

		for x := 0; x < f.width; x++ {
			heat := f.buffer[y*f.width+x]

			// Skip very low heat (natural fade to background)
			if heat < 3 {
				// Flush any pending batch
				if batchChars.Len() > 0 {
					r, g, b := hexToRGB(currentColor)
					fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
					batchChars.Reset()
					currentColor = ""
				}
				output.WriteRune(' ')
				continue
			}

			// Map heat to character (0-36 heat → 5 chars)
			charIndex := heat / 7
			if charIndex >= len(f.chars) {
				charIndex = len(f.chars) - 1
			}
			char := f.chars[charIndex]

			// Map heat to color from palette
			colorIndex := heat * (len(f.palette) - 1) / 36
			if colorIndex >= len(f.palette) {
				colorIndex = len(f.palette) - 1
			}
			color := f.palette[colorIndex]

			// If color changed, flush current batch
			if color != currentColor && batchChars.Len() > 0 {
				r, g, b := hexToRGB(currentColor)
				fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
				batchChars.Reset()
			}

			// Add char to batch
			currentColor = color
			batchChars.WriteRune(char)
		}

		// Flush final batch for this row
		if batchChars.Len() > 0 {
			r, g, b := hexToRGB(currentColor)
			fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, batchChars.String())
		}

		// Add newline at end of row (but not after last row)
		if y < f.height-1 {
			output.WriteRune('\n')
		}
	}

	return output.String()
}
