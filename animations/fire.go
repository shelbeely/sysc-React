package animations

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
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

// spreadFire propagates heat upward with random spread and decay
func (f *FireEffect) spreadFire(from int) {
	// Random horizontal offset (-1 to +2) for natural flame movement
	offset := rand.Intn(4) - 1
	to := from - f.width + offset

	// Bounds check
	if to < 0 || to >= len(f.buffer) {
		return
	}

	// Random decay (0 or 1) - Ly style
	decay := rand.Intn(2)

	// Apply decay
	newHeat := f.buffer[from] - decay
	if newHeat < 0 {
		newHeat = 0
	}

	f.buffer[to] = newHeat
}

// Update advances fire simulation - bottom-to-top heat propagation
func (f *FireEffect) Update() {
	// Randomly re-ignite bottom row
	for x := 0; x < f.width; x++ {
		if rand.Float64() < 0.5 {
			f.buffer[(f.height-1)*f.width+x] = 65
		}
	}

	// Propagate fire from bottom to top
	for y := f.height - 1; y > 0; y-- {
		for x := 0; x < f.width; x++ {
			index := y*f.width + x
			f.spreadFire(index)
		}
	}
}

// Render converts fire to colored block output with batched styling
func (f *FireEffect) Render() string {
	var lines []string

	for y := 0; y < f.height; y++ {
		var line strings.Builder

		// Batch consecutive chars with same color
		var batchChars strings.Builder
		var batchColor string

		for x := 0; x < f.width; x++ {
			intensity := f.buffer[y*f.width+x]

			// Handle low intensity (spaces)
			if intensity < 5 {
				// Flush any pending batch
				if batchChars.Len() > 0 {
					styled := lipgloss.NewStyle().
						Foreground(lipgloss.Color(batchColor)).
						Render(batchChars.String())
					line.WriteString(styled)
					batchChars.Reset()
				}
				line.WriteRune(' ')
				batchColor = ""
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

			// If color changed, flush current batch and start new one
			if color != batchColor && batchChars.Len() > 0 {
				styled := lipgloss.NewStyle().
					Foreground(lipgloss.Color(batchColor)).
					Render(batchChars.String())
				line.WriteString(styled)
				batchChars.Reset()
			}

			// Add char to batch
			batchColor = color
			batchChars.WriteRune(char)
		}

		// Flush final batch
		if batchChars.Len() > 0 {
			styled := lipgloss.NewStyle().
				Foreground(lipgloss.Color(batchColor)).
				Render(batchChars.String())
			line.WriteString(styled)
		}

		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
