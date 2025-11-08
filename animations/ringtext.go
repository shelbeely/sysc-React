package animations

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
)

// RingTextConfig holds the configuration for the RingText effect
type RingTextConfig struct {
	Width              int
	Height             int
	Text               string
	RingColors         []string  // Colors for each ring
	RingGap            float64   // Distance between rings as a percent of smallest dimension
	SpinSpeed          float64   // Speed of rotation (radians per frame)
	DisperseFrames     int       // Frames to transition to rings
	SpinFrames         int       // Frames to spin before returning
	ReturnFrames       int       // Frames to transition back to text
	StaticFrames       int       // Frames to display static text initially
	FinalGradientStops []string  // Gradient for final text state
	FinalGradientSteps int       // Number of gradient steps
}

// RingTextEffect represents the multi-phase ring text animation
type RingTextEffect struct {
	width  int
	height int
	text   string

	// Ring configuration
	ringColors     []string
	ringGap        float64
	spinSpeed      float64
	disperseFrames int
	spinFrames     int
	returnFrames   int
	staticFrames   int

	// Gradient configuration
	finalGradientStops []string
	finalGradientSteps int
	finalGradient      []string

	// Character data
	chars      []RingTextCharacter
	rings      []Ring
	centerX    float64
	centerY    float64
	rng        *rand.Rand
	frameCount int

	// Animation state
	phase string // "static", "disperse_to_rings", "spin", "return_to_text", "hold"
}

// RingTextCharacter represents a single character in the animation
type RingTextCharacter struct {
	original      rune
	x             int     // Original position
	y             int     // Original position
	currentX      float64 // Current animated position
	currentY      float64 // Current animated position
	visible       bool
	currentColor  string
	ringIndex     int     // Which ring this character belongs to
	angleOnRing   float64 // Current angle on the ring (in radians)
	targetAngle   float64 // Target angle for smooth transitions
	rotationSpeed float64 // Individual rotation speed
	clockwise     bool    // Rotation direction
}

// Ring represents a circular ring of positions
type Ring struct {
	radius        float64
	color         string
	rotationSpeed float64
	clockwise     bool
	characterIndices []int // Indices of characters on this ring
}

// NewRingTextEffect creates a new RingText effect
func NewRingTextEffect(config RingTextConfig) *RingTextEffect {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Set defaults
	if config.RingGap == 0 {
		config.RingGap = 0.15
	}
	if config.SpinSpeed == 0 {
		config.SpinSpeed = 0.05
	}
	if config.DisperseFrames == 0 {
		config.DisperseFrames = 100
	}
	if config.SpinFrames == 0 {
		config.SpinFrames = 200
	}
	if config.ReturnFrames == 0 {
		config.ReturnFrames = 100
	}
	if config.StaticFrames == 0 {
		config.StaticFrames = 60
	}
	if config.FinalGradientSteps == 0 {
		config.FinalGradientSteps = 12
	}
	if len(config.RingColors) == 0 {
		config.RingColors = []string{"#bd93f9", "#ff79c6", "#f1fa8c"} // Default: Dracula purple, pink, yellow
	}
	if len(config.FinalGradientStops) == 0 {
		config.FinalGradientStops = []string{"#6272a4", "#bd93f9", "#f8f8f2"}
	}

	effect := &RingTextEffect{
		width:              config.Width,
		height:             config.Height,
		text:               config.Text,
		ringColors:         config.RingColors,
		ringGap:            config.RingGap,
		spinSpeed:          config.SpinSpeed,
		disperseFrames:     config.DisperseFrames,
		spinFrames:         config.SpinFrames,
		returnFrames:       config.ReturnFrames,
		staticFrames:       config.StaticFrames,
		finalGradientStops: config.FinalGradientStops,
		finalGradientSteps: config.FinalGradientSteps,
		rng:                rng,
		phase:              "static",
		frameCount:         0,
	}

	effect.init()
	return effect
}

// init initializes the effect
func (e *RingTextEffect) init() {
	e.centerX = float64(e.width) / 2
	e.centerY = float64(e.height) / 2

	// Create gradient for final state
	e.finalGradient = e.createGradient(e.finalGradientStops, e.finalGradientSteps)

	// Parse text and create characters
	e.parseText()

	// Create rings and assign characters
	e.createRings()
}

// parseText converts the text into positioned characters
func (e *RingTextEffect) parseText() {
	lines := strings.Split(e.text, "\n")
	totalLines := len(lines)

	// Calculate starting Y position to center text vertically
	startY := (e.height - totalLines) / 2

	e.chars = make([]RingTextCharacter, 0)

	for lineIdx, line := range lines {
		lineRunes := []rune(line)
		lineLen := len(lineRunes)

		// Calculate starting X position to center line horizontally
		startX := (e.width - lineLen) / 2

		for charIdx, char := range lineRunes {
			if char == ' ' || char == '\n' {
				continue // Skip spaces and newlines
			}

			x := startX + charIdx
			y := startY + lineIdx

			character := RingTextCharacter{
				original: char,
				x:        x,
				y:        y,
				currentX: float64(x),
				currentY: float64(y),
				visible:  true,
				currentColor: e.finalGradient[0], // Start with first gradient color
			}

			e.chars = append(e.chars, character)
		}
	}
}

// createRings creates concentric rings and assigns characters to them
func (e *RingTextEffect) createRings() {
	if len(e.chars) == 0 {
		return
	}

	// Calculate maximum radius based on smallest dimension
	smallestDim := float64(e.width)
	if float64(e.height) < smallestDim {
		smallestDim = float64(e.height)
	}

	ringGapPixels := smallestDim * e.ringGap
	maxRadius := smallestDim / 2

	// Create rings
	e.rings = make([]Ring, 0)
	for radius := ringGapPixels; radius < maxRadius; radius += ringGapPixels {
		colorIndex := len(e.rings) % len(e.ringColors)
		clockwise := len(e.rings)%2 == 0

		ring := Ring{
			radius:           radius,
			color:            e.ringColors[colorIndex],
			rotationSpeed:    e.spinSpeed * (0.5 + e.rng.Float64()*0.5), // Randomize speed slightly
			clockwise:        clockwise,
			characterIndices: make([]int, 0),
		}

		e.rings = append(e.rings, ring)
	}

	// Assign characters to rings evenly
	if len(e.rings) > 0 {
		for i := range e.chars {
			ringIndex := i % len(e.rings)
			e.chars[i].ringIndex = ringIndex
			e.chars[i].clockwise = e.rings[ringIndex].clockwise
			e.chars[i].rotationSpeed = e.rings[ringIndex].rotationSpeed

			// Calculate initial angle on ring based on character's original position
			dx := float64(e.chars[i].x) - e.centerX
			dy := float64(e.chars[i].y) - e.centerY
			e.chars[i].angleOnRing = math.Atan2(dy, dx)
			e.chars[i].targetAngle = e.chars[i].angleOnRing

			e.rings[ringIndex].characterIndices = append(e.rings[ringIndex].characterIndices, i)
		}
	}
}

// Update advances the animation by one frame
func (e *RingTextEffect) Update() {
	e.frameCount++

	switch e.phase {
	case "static":
		if e.frameCount >= e.staticFrames {
			e.phase = "disperse_to_rings"
			e.frameCount = 0
		}

	case "disperse_to_rings":
		progress := float64(e.frameCount) / float64(e.disperseFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		// Ease-in-out function for smooth transition
		easedProgress := e.easeInOutCubic(progress)

		for i := range e.chars {
			ring := &e.rings[e.chars[i].ringIndex]

			// Calculate target position on ring
			targetX := e.centerX + ring.radius*math.Cos(e.chars[i].angleOnRing)
			targetY := e.centerY + ring.radius*math.Sin(e.chars[i].angleOnRing)

			// Interpolate position
			e.chars[i].currentX = float64(e.chars[i].x) + (targetX-float64(e.chars[i].x))*easedProgress
			e.chars[i].currentY = float64(e.chars[i].y) + (targetY-float64(e.chars[i].y))*easedProgress

			// Interpolate color to ring color
			e.chars[i].currentColor = e.interpolateColor(e.finalGradient[0], ring.color, easedProgress)
		}

		if e.frameCount >= e.disperseFrames {
			e.phase = "spin"
			e.frameCount = 0
		}

	case "spin":
		// Rotate characters around their rings
		for i := range e.chars {
			if e.chars[i].clockwise {
				e.chars[i].angleOnRing += e.chars[i].rotationSpeed
			} else {
				e.chars[i].angleOnRing -= e.chars[i].rotationSpeed
			}

			// Normalize angle to [0, 2π]
			for e.chars[i].angleOnRing > 2*math.Pi {
				e.chars[i].angleOnRing -= 2 * math.Pi
			}
			for e.chars[i].angleOnRing < 0 {
				e.chars[i].angleOnRing += 2 * math.Pi
			}

			ring := &e.rings[e.chars[i].ringIndex]
			e.chars[i].currentX = e.centerX + ring.radius*math.Cos(e.chars[i].angleOnRing)
			e.chars[i].currentY = e.centerY + ring.radius*math.Sin(e.chars[i].angleOnRing)
		}

		if e.frameCount >= e.spinFrames {
			e.phase = "return_to_text"
			e.frameCount = 0
		}

	case "return_to_text":
		progress := float64(e.frameCount) / float64(e.returnFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		// Ease-in-out function for smooth transition
		easedProgress := e.easeInOutCubic(progress)

		for i := range e.chars {
			ring := &e.rings[e.chars[i].ringIndex]

			// Calculate current ring position
			ringX := e.centerX + ring.radius*math.Cos(e.chars[i].angleOnRing)
			ringY := e.centerY + ring.radius*math.Sin(e.chars[i].angleOnRing)

			// Interpolate back to original position
			e.chars[i].currentX = ringX + (float64(e.chars[i].x)-ringX)*easedProgress
			e.chars[i].currentY = ringY + (float64(e.chars[i].y)-ringY)*easedProgress

			// Interpolate color back to final gradient color
			gradientIndex := int(easedProgress * float64(len(e.finalGradient)-1))
			if gradientIndex >= len(e.finalGradient) {
				gradientIndex = len(e.finalGradient) - 1
			}
			e.chars[i].currentColor = e.finalGradient[gradientIndex]
		}

		if e.frameCount >= e.returnFrames {
			e.phase = "hold"
			e.frameCount = 0
		}

	case "hold":
		// Hold the final state for a bit before looping
		if e.frameCount >= 60 {
			e.Reset()
		}
	}
}

// Render returns the current frame as a colored string
func (e *RingTextEffect) Render() string {
	// Create a 2D buffer for the screen
	buffer := make([][]rune, e.height)
	colors := make([][]string, e.height)
	for i := range buffer {
		buffer[i] = make([]rune, e.width)
		colors[i] = make([]string, e.width)
		for j := range buffer[i] {
			buffer[i][j] = ' '
			colors[i][j] = ""
		}
	}

	// Draw characters
	for _, char := range e.chars {
		if !char.visible {
			continue
		}

		x := int(math.Round(char.currentX))
		y := int(math.Round(char.currentY))

		// Bounds check
		if x >= 0 && x < e.width && y >= 0 && y < e.height {
			buffer[y][x] = char.original
			colors[y][x] = char.currentColor
		}
	}

	// Build output string
	var output strings.Builder
	for y := 0; y < e.height; y++ {
		for x := 0; x < e.width; x++ {
			char := buffer[y][x]
			color := colors[y][x]

			if color != "" && char != ' ' {
				styled := lipgloss.NewStyle().
					Foreground(lipgloss.Color(color)).
					Render(string(char))
				output.WriteString(styled)
			} else {
				output.WriteRune(char)
			}
		}
		if y < e.height-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}

// Reset restarts the animation
func (e *RingTextEffect) Reset() {
	e.phase = "static"
	e.frameCount = 0

	// Reset character positions
	for i := range e.chars {
		e.chars[i].currentX = float64(e.chars[i].x)
		e.chars[i].currentY = float64(e.chars[i].y)
		e.chars[i].currentColor = e.finalGradient[0]

		// Reset angle
		dx := float64(e.chars[i].x) - e.centerX
		dy := float64(e.chars[i].y) - e.centerY
		e.chars[i].angleOnRing = math.Atan2(dy, dx)
	}
}

// createGradient creates a gradient between color stops
func (e *RingTextEffect) createGradient(stops []string, steps int) []string {
	if len(stops) == 0 {
		return []string{"#ffffff"}
	}
	if len(stops) == 1 {
		return []string{stops[0]}
	}

	gradient := make([]string, 0)
	stepsPerSegment := steps / (len(stops) - 1)

	for i := 0; i < len(stops)-1; i++ {
		startColor := parseHexColor(stops[i])
		endColor := parseHexColor(stops[i+1])

		for j := 0; j < stepsPerSegment; j++ {
			t := float64(j) / float64(stepsPerSegment)
			r := uint8(float64(startColor[0]) + (float64(endColor[0])-float64(startColor[0]))*t)
			g := uint8(float64(startColor[1]) + (float64(endColor[1])-float64(startColor[1]))*t)
			b := uint8(float64(startColor[2]) + (float64(endColor[2])-float64(startColor[2]))*t)
			gradient = append(gradient, formatHexColor([3]uint8{r, g, b}))
		}
	}

	// Add final color
	gradient = append(gradient, stops[len(stops)-1])
	return gradient
}

// interpolateColor interpolates between two hex colors
func (e *RingTextEffect) interpolateColor(color1, color2 string, t float64) string {
	c1 := parseHexColor(color1)
	c2 := parseHexColor(color2)

	r := uint8(float64(c1[0]) + (float64(c2[0])-float64(c1[0]))*t)
	g := uint8(float64(c1[1]) + (float64(c2[1])-float64(c1[1]))*t)
	b := uint8(float64(c1[2]) + (float64(c2[2])-float64(c1[2]))*t)

	return formatHexColor([3]uint8{r, g, b})
}

// easeInOutCubic applies an ease-in-out cubic easing function
func (e *RingTextEffect) easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}
