package animations

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// FireEffect implements PSX DOOM-style fire algorithm
type FireEffect struct {
	width   int      // Terminal width
	height  int      // Terminal height
	buffer  []int    // Heat values (0-35), size = width * height
	palette []string // Hex color codes from theme
	chars   []rune   // Fire characters for density
}

const (
	fireSteps     = 36 // Fire intensity levels (0-35)
	fireSpread    = 3  // Maximum horizontal spread distance
	fireDecayRate = 7  // Decay rate out of 10 (higher = faster decay)
)

// NewFireEffect creates a new fire effect with given dimensions and theme palette
func NewFireEffect(width, height int, palette []string) *FireEffect {
	f := &FireEffect{
		width:   width,
		height:  height,
		palette: palette,
		// Block characters with increasing density for fire intensity
		chars:   []rune{'░', '▒', '▓', '█', '▓', '█', '█', '█'},
	}
	f.init()
	return f
}

// Initialize fire buffer with bottom row as heat source
func (f *FireEffect) init() {
	f.buffer = make([]int, f.width*f.height)

	// Set bottom row to maximum heat with some randomness (fire source)
	for i := 0; i < f.width; i++ {
		f.buffer[(f.height-1)*f.width+i] = fireSteps - 1 - rand.Intn(4)
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

// spreadFire propagates heat upward with doom-style random decay
func (f *FireEffect) spreadFire(from int) {
	// Random horizontal spread (-fireSpread to +fireSpread)
	randSpread := rand.Intn(fireSpread*2+1) - fireSpread
	to := from - f.width + randSpread

	// Bounds check
	if to < 0 || to >= len(f.buffer) {
		return
	}

	// Random decay: fireDecayRate out of 10 chance to decay
	// This creates the characteristic flickering doom fire effect
	decay := 0
	if rand.Intn(10) < fireDecayRate {
		decay = 1
	}

	newHeat := f.buffer[from] - decay
	if newHeat < 0 {
		newHeat = 0
	}

	f.buffer[to] = newHeat
}

// Update advances the fire simulation by one frame
func (f *FireEffect) Update() {
	// Randomly ignite bottom row (fire source) with varying intensity
	for i := 0; i < f.width/8; i++ {
		pos := rand.Intn(f.width)
		f.buffer[(f.height-1)*f.width+pos] = fireSteps - 1 - rand.Intn(8)
	}

	// Process all pixels from bottom to top
	// (Fire spreads upward, must process bottom row first)
	for y := f.height - 1; y > 0; y-- {
		for x := 0; x < f.width; x++ {
			index := y*f.width + x
			f.spreadFire(index)
		}
	}
}

// Render converts the fire buffer to colored text output
func (f *FireEffect) Render() string {
	var lines []string

	for y := 0; y < f.height; y++ {
		var line strings.Builder
		for x := 0; x < f.width; x++ {
			heat := f.buffer[y*f.width+x]

			// Skip zero/very low heat (natural fade to background)
			if heat < 2 {
				line.WriteRune(' ')
				continue
			}

			// Map heat to character (0-35 heat → block density)
			charIndex := (heat * (len(f.chars) - 1)) / (fireSteps - 1)
			if charIndex >= len(f.chars) {
				charIndex = len(f.chars) - 1
			}
			char := f.chars[charIndex]

			// Skip styling if character is a space (save massive amounts of ANSI codes)
			if char == ' ' {
				line.WriteRune(' ')
				continue
			}

			// Map heat to color from palette
			colorIndex := (heat * (len(f.palette) - 1)) / (fireSteps - 1)
			if colorIndex >= len(f.palette) {
				colorIndex = len(f.palette) - 1
			}
			colorHex := f.palette[colorIndex]

			// Render colored character
			styled := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorHex)).
				Render(string(char))
			line.WriteString(styled)
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
