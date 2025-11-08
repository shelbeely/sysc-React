package animations

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
)

// BlackholeConfig holds the configuration for the Blackhole effect
type BlackholeConfig struct {
	Width              int
	Height             int
	Text               string
	BlackholeColor     string   // Border color for singularity
	StarColors         []string // Colors for post-explosion stars
	FinalGradientStops []string // Gradient for final text state
	FinalGradientSteps int      // Number of gradient steps
	FinalGradientDir   GradientDirection
	StaticGradientStops []string // Gradient for static ASCII
	StaticGradientDir   GradientDirection
	FormingFrames      int // Frames for border formation
	ConsumingFrames    int // Frames for consumption
	CollapsingFrames   int // Frames for border collapse
	ExplodingFrames    int // Frames for explosion scatter
	ReturningFrames    int // Frames for return to text
	StaticFrames       int // Frames to display static text initially
}

// BlackholeEffect represents the multi-phase blackhole animation
type BlackholeEffect struct {
	width  int
	height int
	text   string

	// Blackhole configuration
	blackholeColor     string
	starColors         []string
	finalGradientStops []string
	finalGradientSteps int
	finalGradientDir   GradientDirection
	staticGradientStops []string
	staticGradientDir   GradientDirection
	formingFrames      int
	consumingFrames    int
	collapsingFrames   int
	explodingFrames    int
	returningFrames    int
	staticFrames       int

	// Gradients
	finalGradient  []string
	staticGradient []string
	starGradient   []string

	// Character data
	chars          []BlackholeCharacter
	borderChars    []BorderCharacter
	centerX        float64
	centerY        float64
	blackholeRadius float64
	rng            *rand.Rand
	frameCount     int

	// Animation state
	phase          string // "static", "forming", "consuming", "collapsing", "exploding", "returning", "hold"
	consumeCounter int    // Track consumption progress
}

// BlackholeCharacter represents a single character in the animation
type BlackholeCharacter struct {
	original     rune
	x            int     // Original position
	y            int     // Original position
	currentX     float64 // Current animated position
	currentY     float64 // Current animated position
	scatterX     float64 // Scatter position for explosion
	scatterY     float64 // Scatter position for explosion
	visible      bool
	currentColor string
	consumed     bool    // Has been consumed by blackhole
	consumeOrder int     // Order in which character is consumed
	scatterAngle float64 // Direction for explosion scatter
	scatterDist  float64 // Distance for explosion scatter
}

// BorderCharacter represents a character on the blackhole border
type BorderCharacter struct {
	angle        float64
	currentX     float64
	currentY     float64
	symbol       rune
	currentColor string
	visible      bool
}

var unstableSymbols = []rune{'◦', '◎', '◉', '●'}

// NewBlackholeEffect creates a new Blackhole effect
func NewBlackholeEffect(config BlackholeConfig) *BlackholeEffect {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Set defaults
	if config.BlackholeColor == "" {
		config.BlackholeColor = "#ffffff"
	}
	if len(config.StarColors) == 0 {
		config.StarColors = []string{"#ffffff", "#ffd700", "#ff6b6b", "#4ecdc4", "#95e1d3", "#f38181"}
	}
	if len(config.FinalGradientStops) == 0 {
		config.FinalGradientStops = []string{"#6272a4", "#bd93f9", "#f8f8f2"}
	}
	if config.FinalGradientSteps == 0 {
		config.FinalGradientSteps = 12
	}
	if len(config.StaticGradientStops) == 0 {
		config.StaticGradientStops = config.StarColors
	}
	if config.FormingFrames == 0 {
		config.FormingFrames = 100
	}
	if config.ConsumingFrames == 0 {
		config.ConsumingFrames = 150
	}
	if config.CollapsingFrames == 0 {
		config.CollapsingFrames = 50
	}
	if config.ExplodingFrames == 0 {
		config.ExplodingFrames = 100
	}
	if config.ReturningFrames == 0 {
		config.ReturningFrames = 120
	}
	if config.StaticFrames == 0 {
		config.StaticFrames = 100
	}

	effect := &BlackholeEffect{
		width:               config.Width,
		height:              config.Height,
		text:                config.Text,
		blackholeColor:      config.BlackholeColor,
		starColors:          config.StarColors,
		finalGradientStops:  config.FinalGradientStops,
		finalGradientSteps:  config.FinalGradientSteps,
		finalGradientDir:    config.FinalGradientDir,
		staticGradientStops: config.StaticGradientStops,
		staticGradientDir:   config.StaticGradientDir,
		formingFrames:       config.FormingFrames,
		consumingFrames:     config.ConsumingFrames,
		collapsingFrames:    config.CollapsingFrames,
		explodingFrames:     config.ExplodingFrames,
		returningFrames:     config.ReturningFrames,
		staticFrames:        config.StaticFrames,
		rng:                 rng,
		phase:               "static",
		frameCount:          0,
		consumeCounter:      0,
	}

	effect.init()
	return effect
}

// init initializes the effect
func (e *BlackholeEffect) init() {
	e.centerX = float64(e.width) / 2
	e.centerY = float64(e.height) / 2

	// Calculate blackhole radius (30% of smallest dimension, minimum 3)
	smallestDim := float64(e.width)
	if float64(e.height) < smallestDim {
		smallestDim = float64(e.height)
	}
	e.blackholeRadius = math.Max(smallestDim*0.3, 3)

	// Create gradients
	e.finalGradient = e.createGradient(e.finalGradientStops, e.finalGradientSteps)
	e.staticGradient = e.createGradient(e.staticGradientStops, 100)
	e.starGradient = e.createGradient(e.starColors, 100)

	// Parse text and create characters
	e.parseText()

	// Create border characters
	e.createBorder()

	// Apply initial static gradient
	e.applyStaticGradient()

	// Generate scatter positions for explosion
	e.generateScatterPositions()
}

// parseText converts the text into positioned characters
func (e *BlackholeEffect) parseText() {
	lines := strings.Split(e.text, "\n")
	totalLines := len(lines)

	startY := (e.height - totalLines) / 2

	e.chars = make([]BlackholeCharacter, 0)

	for lineIdx, line := range lines {
		lineRunes := []rune(line)
		lineLen := len(lineRunes)
		startX := (e.width - lineLen) / 2

		for charIdx, char := range lineRunes {
			if char == ' ' || char == '\n' {
				continue
			}

			x := startX + charIdx
			y := startY + lineIdx

			character := BlackholeCharacter{
				original:     char,
				x:            x,
				y:            y,
				currentX:     float64(x),
				currentY:     float64(y),
				visible:      true,
				currentColor: e.staticGradient[0],
				consumed:     false,
				consumeOrder: -1,
			}

			e.chars = append(e.chars, character)
		}
	}

	// Assign random consumption order
	indices := make([]int, len(e.chars))
	for i := range indices {
		indices[i] = i
	}
	// Fisher-Yates shuffle
	for i := len(indices) - 1; i > 0; i-- {
		j := e.rng.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}
	for order, idx := range indices {
		e.chars[idx].consumeOrder = order
	}
}

// createBorder creates the circular border around the blackhole
func (e *BlackholeEffect) createBorder() {
	// Create approximately 50 border characters around the circle
	numBorderChars := 50
	e.borderChars = make([]BorderCharacter, numBorderChars)

	for i := range e.borderChars {
		angle := (float64(i) / float64(numBorderChars)) * 2 * math.Pi
		e.borderChars[i] = BorderCharacter{
			angle:        angle,
			currentX:     e.centerX + e.blackholeRadius*math.Cos(angle),
			currentY:     e.centerY + e.blackholeRadius*math.Sin(angle),
			symbol:       '●',
			currentColor: e.blackholeColor,
			visible:      false,
		}
	}
}

// generateScatterPositions creates random scatter positions for explosion
func (e *BlackholeEffect) generateScatterPositions() {
	for i := range e.chars {
		// Random angle for scatter direction
		e.chars[i].scatterAngle = e.rng.Float64() * 2 * math.Pi

		// Random distance (50-150% of blackhole radius)
		e.chars[i].scatterDist = e.blackholeRadius * (0.5 + e.rng.Float64())

		// Calculate scatter position
		e.chars[i].scatterX = e.centerX + math.Cos(e.chars[i].scatterAngle)*e.chars[i].scatterDist
		e.chars[i].scatterY = e.centerY + math.Sin(e.chars[i].scatterAngle)*e.chars[i].scatterDist

		// Clamp to canvas
		if e.chars[i].scatterX < 0 {
			e.chars[i].scatterX = 0
		}
		if e.chars[i].scatterX >= float64(e.width) {
			e.chars[i].scatterX = float64(e.width - 1)
		}
		if e.chars[i].scatterY < 0 {
			e.chars[i].scatterY = 0
		}
		if e.chars[i].scatterY >= float64(e.height) {
			e.chars[i].scatterY = float64(e.height - 1)
		}
	}
}

// applyStaticGradient applies gradient to static ASCII (same as ringtext)
func (e *BlackholeEffect) applyStaticGradient() {
	if len(e.chars) == 0 || len(e.staticGradient) == 0 {
		return
	}

	minX, maxX := e.width, 0
	minY, maxY := e.height, 0

	for i := range e.chars {
		if e.chars[i].x < minX {
			minX = e.chars[i].x
		}
		if e.chars[i].x > maxX {
			maxX = e.chars[i].x
		}
		if e.chars[i].y < minY {
			minY = e.chars[i].y
		}
		if e.chars[i].y > maxY {
			maxY = e.chars[i].y
		}
	}

	textWidth := float64(maxX - minX)
	textHeight := float64(maxY - minY)
	if textWidth == 0 {
		textWidth = 1
	}
	if textHeight == 0 {
		textHeight = 1
	}

	for i := range e.chars {
		var gradientPos float64

		switch e.staticGradientDir {
		case GradientHorizontal:
			gradientPos = float64(e.chars[i].x-minX) / textWidth
		case GradientVertical:
			gradientPos = float64(e.chars[i].y-minY) / textHeight
		case GradientDiagonal:
			xPos := float64(e.chars[i].x-minX) / textWidth
			yPos := float64(e.chars[i].y-minY) / textHeight
			gradientPos = (xPos + yPos) / 2.0
		case GradientRadial:
			dx := float64(e.chars[i].x) - e.centerX
			dy := float64(e.chars[i].y) - e.centerY
			maxDist := math.Sqrt(textWidth*textWidth+textHeight*textHeight) / 2.0
			dist := math.Sqrt(dx*dx + dy*dy)
			gradientPos = math.Min(dist/maxDist, 1.0)
		default:
			gradientPos = 0
		}

		gradientPos = math.Max(0, math.Min(1, gradientPos))
		gradientIndex := int(gradientPos * float64(len(e.staticGradient)-1))
		if gradientIndex >= len(e.staticGradient) {
			gradientIndex = len(e.staticGradient) - 1
		}

		e.chars[i].currentColor = e.staticGradient[gradientIndex]
	}
}

// Update advances the animation by one frame
func (e *BlackholeEffect) Update() {
	e.frameCount++

	switch e.phase {
	case "static":
		if e.frameCount >= e.staticFrames {
			e.phase = "forming"
			e.frameCount = 0
		}

	case "forming":
		progress := float64(e.frameCount) / float64(e.formingFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		// Gradually show border characters
		visibleCount := int(progress * float64(len(e.borderChars)))
		for i := 0; i < visibleCount && i < len(e.borderChars); i++ {
			e.borderChars[i].visible = true
		}

		if e.frameCount >= e.formingFrames {
			e.phase = "consuming"
			e.frameCount = 0
			e.consumeCounter = 0
		}

	case "consuming":
		progress := float64(e.frameCount) / float64(e.consumingFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		// Consume characters gradually (1-5 per frame based on progress)
		charsPerFrame := 1 + int(progress*4)
		for i := 0; i < charsPerFrame; i++ {
			if e.consumeCounter < len(e.chars) {
				// Find next character to consume
				for j := range e.chars {
					if e.chars[j].consumeOrder == e.consumeCounter && !e.chars[j].consumed {
						e.chars[j].consumed = true
						break
					}
				}
				e.consumeCounter++
			}
		}

		// Move consumed characters toward center with exponential easing
		for i := range e.chars {
			if e.chars[i].consumed {
				// Exponential ease toward center (gravity effect)
				easedProgress := e.easeInExpo(progress)

				// Bézier curve toward center
				startX := float64(e.chars[i].x)
				startY := float64(e.chars[i].y)

				// Control point for curve (offset perpendicular to direction)
				dx := e.centerX - startX
				dy := e.centerY - startY
				controlX := startX + dx*0.5 + dy*0.3
				controlY := startY + dy*0.5 - dx*0.3

				// Quadratic Bézier curve
				t := easedProgress
				e.chars[i].currentX = (1-t)*(1-t)*startX + 2*(1-t)*t*controlX + t*t*e.centerX
				e.chars[i].currentY = (1-t)*(1-t)*startY + 2*(1-t)*t*controlY + t*t*e.centerY

				// Fade to black as approaching center
				dist := math.Sqrt(math.Pow(e.chars[i].currentX-e.centerX, 2) + math.Pow(e.chars[i].currentY-e.centerY, 2))
				brightness := dist / e.blackholeRadius
				if brightness < 0.3 {
					e.chars[i].visible = false
				}
			}
		}

		if e.frameCount >= e.consumingFrames {
			e.phase = "collapsing"
			e.frameCount = 0
		}

	case "collapsing":
		progress := float64(e.frameCount) / float64(e.collapsingFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		// Contract border toward center
		currentRadius := e.blackholeRadius * (1.0 - progress)

		for i := range e.borderChars {
			e.borderChars[i].currentX = e.centerX + currentRadius*math.Cos(e.borderChars[i].angle)
			e.borderChars[i].currentY = e.centerY + currentRadius*math.Sin(e.borderChars[i].angle)

			// Random unstable symbols
			if e.rng.Float64() < 0.1 {
				e.borderChars[i].symbol = unstableSymbols[e.rng.Intn(len(unstableSymbols))]
			}

			// Random colors from star palette
			if e.rng.Float64() < 0.05 {
				e.borderChars[i].currentColor = e.starColors[e.rng.Intn(len(e.starColors))]
			}
		}

		if e.frameCount >= e.collapsingFrames {
			e.phase = "exploding"
			e.frameCount = 0
			// Hide border
			for i := range e.borderChars {
				e.borderChars[i].visible = false
			}
			// Reset character visibility
			for i := range e.chars {
				e.chars[i].visible = true
			}
		}

	case "exploding":
		progress := float64(e.frameCount) / float64(e.explodingFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		easedProgress := e.easeOutExpo(progress)

		for i := range e.chars {
			// Scatter from center to scatter position
			e.chars[i].currentX = e.centerX + (e.chars[i].scatterX-e.centerX)*easedProgress
			e.chars[i].currentY = e.centerY + (e.chars[i].scatterY-e.centerY)*easedProgress

			// Cycle through star colors
			colorIndex := int((progress + float64(i)*0.1) * float64(len(e.starGradient)))
			colorIndex = colorIndex % len(e.starGradient)
			e.chars[i].currentColor = e.starGradient[colorIndex]
		}

		if e.frameCount >= e.explodingFrames {
			e.phase = "returning"
			e.frameCount = 0
		}

	case "returning":
		progress := float64(e.frameCount) / float64(e.returningFrames)
		if progress > 1.0 {
			progress = 1.0
		}

		easedProgress := e.easeInOutCubic(progress)

		for i := range e.chars {
			// Return from scatter position to original
			e.chars[i].currentX = e.chars[i].scatterX + (float64(e.chars[i].x)-e.chars[i].scatterX)*easedProgress
			e.chars[i].currentY = e.chars[i].scatterY + (float64(e.chars[i].y)-e.chars[i].scatterY)*easedProgress

			// Transition to final gradient color
			gradientIndex := int(easedProgress * float64(len(e.finalGradient)-1))
			if gradientIndex >= len(e.finalGradient) {
				gradientIndex = len(e.finalGradient) - 1
			}
			e.chars[i].currentColor = e.finalGradient[gradientIndex]
		}

		if e.frameCount >= e.returningFrames {
			e.phase = "hold"
			e.frameCount = 0
		}

	case "hold":
		if e.frameCount >= 60 {
			e.Reset()
		}
	}
}

// Render returns the current frame as a colored string
func (e *BlackholeEffect) Render() string {
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

		if x >= 0 && x < e.width && y >= 0 && y < e.height {
			buffer[y][x] = char.original
			colors[y][x] = char.currentColor
		}
	}

	// Draw border
	for _, borderChar := range e.borderChars {
		if !borderChar.visible {
			continue
		}

		x := int(math.Round(borderChar.currentX))
		y := int(math.Round(borderChar.currentY))

		if x >= 0 && x < e.width && y >= 0 && y < e.height {
			buffer[y][x] = borderChar.symbol
			colors[y][x] = borderChar.currentColor
		}
	}

	// Build output
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
func (e *BlackholeEffect) Reset() {
	e.phase = "static"
	e.frameCount = 0
	e.consumeCounter = 0

	// Reset characters
	for i := range e.chars {
		e.chars[i].currentX = float64(e.chars[i].x)
		e.chars[i].currentY = float64(e.chars[i].y)
		e.chars[i].visible = true
		e.chars[i].consumed = false
	}

	// Reset border
	for i := range e.borderChars {
		e.borderChars[i].visible = false
		e.borderChars[i].symbol = '●'
		e.borderChars[i].currentColor = e.blackholeColor
		e.borderChars[i].currentX = e.centerX + e.blackholeRadius*math.Cos(e.borderChars[i].angle)
		e.borderChars[i].currentY = e.centerY + e.blackholeRadius*math.Sin(e.borderChars[i].angle)
	}

	// Reapply static gradient
	e.applyStaticGradient()

	// Regenerate scatter positions
	e.generateScatterPositions()
}

// createGradient creates a gradient between color stops
func (e *BlackholeEffect) createGradient(stops []string, steps int) []string {
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

	gradient = append(gradient, stops[len(stops)-1])
	return gradient
}

// Easing functions
func (e *BlackholeEffect) easeInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*(t-1))
}

func (e *BlackholeEffect) easeOutExpo(t float64) float64 {
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

func (e *BlackholeEffect) easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}
