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
		// Enhanced 8-character gradient for smoother fire rendering
		chars:  []rune{' ', '░', '▒', '▒', '▓', '▓', '█', '█'},
		buffer: make([]int, width*height),
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

// Update advances fire simulation - hybrid SCRFIRE averaging + Ly decay
func (f *FireEffect) Update() {
	// Randomly ignite bottom row (SCRFIRE style)
	for x := 0; x < f.width; x++ {
		if rand.Float64() < 0.5 {
			f.buffer[(f.height-1)*f.width+x] = 65
		}
	}

	// Diffuse fire using averaging (SCRFIRE) with decay
	for y := 0; y < f.height-1; y++ {
		for x := 0; x < f.width; x++ {
			i := y*f.width + x

			// Get neighbor values for averaging
			current := f.buffer[i]
			right := 0
			below := 0
			diagBelow := 0

			if x+1 < f.width {
				right = f.buffer[i+1]
			}
			if y+1 < f.height {
				below = f.buffer[i+f.width]
			}
			if x+1 < f.width && y+1 < f.height {
				diagBelow = f.buffer[i+f.width+1]
			}

			// Average with neighbors and decay (SCRFIRE diffusion)
			avg := (current + right + below + diagBelow) / 4

			// Apply probabilistic decay (Ly style)
			if rand.Float64() < 0.2 && avg > 0 {
				avg--
			}

			f.buffer[i] = avg
		}
	}
}

// Render converts fire to colored block output
func (f *FireEffect) Render() string {
	var lines []string

	for y := 0; y < f.height; y++ {
		var line strings.Builder

		for x := 0; x < f.width; x++ {
			intensity := f.buffer[y*f.width+x]

			// Skip zero intensity
			if intensity < 5 {
				line.WriteRune(' ')
				continue
			}

			// Map intensity to block character (0-65 → 8 blocks)
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

			// Render colored block
			styled := lipgloss.NewStyle().
				Foreground(lipgloss.Color(color)).
				Render(string(char))
			line.WriteString(styled)
		}

		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
