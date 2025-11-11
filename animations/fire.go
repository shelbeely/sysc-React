package animations

import (
	"fmt"
	"math/rand"
	"strings"
)

// FireEffect - hybrid doom fire combining Ly, SCRFIRE, and optimization
type FireEffect struct {
	width   int
	height  int
	buffer  []int    // Intensity values (0-65)
	palette []string
	chars   []rune
}

// NewFireEffect creates a new fire effect
func NewFireEffect(width, height int, palette []string) *FireEffect {
	f := &FireEffect{
		width:   width,
		height:  height,
		palette: palette,
		chars:   []rune{'░', '▒', '▓', '█'},
		buffer:  make([]int, width*height),
	}

	// Initialize fire source at bottom
	for x := 0; x < width; x++ {
		f.buffer[(height-1)*width+x] = 65
	}

	return f
}

// UpdatePalette changes the color palette
func (f *FireEffect) UpdatePalette(palette []string) {
	f.palette = palette
}

// Resize reinitializes with new dimensions
func (f *FireEffect) Resize(width, height int) {
	f.width = width
	f.height = height
	f.buffer = make([]int, width*height)

	for x := 0; x < width; x++ {
		f.buffer[(height-1)*width+x] = 65
	}
}

// Update advances fire simulation - top-to-bottom pulling heat from below
func (f *FireEffect) Update() {
	// Randomly re-ignite bottom row
	for x := 0; x < f.width; x++ {
		if rand.Float64() < 0.5 {
			f.buffer[(f.height-1)*f.width+x] = 65
		}
	}

	// Process top-to-bottom, each pixel PULLS heat from below
	for y := 0; y < f.height-1; y++ {
		for x := 0; x < f.width; x++ {
			// Random horizontal offset to pull from
			offset := rand.Intn(3) - 1 // -1, 0, or 1
			sourceX := x + offset

			// Bounds check
			if sourceX < 0 || sourceX >= f.width {
				sourceX = x // Fall back to directly below
			}

			// Pull heat from pixel below with decay
			sourceIndex := (y + 1) * f.width + sourceX
			destIndex := y * f.width + x

			decay := rand.Intn(2)
			newHeat := f.buffer[sourceIndex] - decay
			if newHeat < 0 {
				newHeat = 0
			}

			f.buffer[destIndex] = newHeat
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
			intensity := f.buffer[y*f.width+x]

			// Handle low intensity (spaces)
			if intensity < 5 {
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

			// Map intensity to block character (0-65 → 4 blocks)
			charIndex := (intensity * (len(f.chars) - 1)) / 65
			if charIndex >= len(f.chars) {
				charIndex = len(f.chars) - 1
			}
			char := f.chars[charIndex]

			// Map intensity to color
			colorIndex := (intensity * (len(f.palette) - 1)) / 65
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
