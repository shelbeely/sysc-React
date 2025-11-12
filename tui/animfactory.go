package tui

import (
	"os"
	"time"

	"github.com/Nomadcxx/sysc-Go/animations"
)

// AnimationWrapper wraps any animation type to provide a common interface
type AnimationWrapper struct {
	render func() string
	update func()
}

func (a *AnimationWrapper) Update() {
	if a.update != nil {
		a.update()
	}
}

func (a *AnimationWrapper) Render() string {
	if a.render != nil {
		return a.render()
	}
	return ""
}

func (a *AnimationWrapper) Reset() {
	// Not implemented for most animations
}

// createAnimation creates an animation instance based on the selected type and settings
// Returns nil if the animation requires user interaction (editors) or isn't supported yet
func (m *Model) createAnimation() animations.Animation {
	animName := m.animations[m.selectedAnimation]
	themeName := m.themes[m.selectedTheme]
	fileName := m.files[m.selectedFile]

	width := m.canvasHeight * 2 // Rough estimate for width based on height
	if width > m.width-10 {
		width = m.width - 10
	}
	height := m.canvasHeight

	// Handle editor modes
	if fileName == "BIT Text Editor" {
		m.bitEditorMode = true
		// Ensure font is loaded when entering BIT editor
		if m.bitCurrentFont == nil && len(m.bitFonts) > 0 {
			fontPath, err := FindFontPath(m.bitFonts[m.bitSelectedFont])
			if err == nil {
				font, err := LoadBitFont(fontPath)
				if err == nil {
					m.bitCurrentFont = font
				}
			}
		}
		m.bitTextInput.Focus()
		return nil
	}
	if fileName == "Custom text" {
		m.editorMode = true
		return nil
	}

	// Create animation based on type (only simple constructors for now)
	switch animName {
	case "fire":
		palette := animations.GetFirePalette(themeName)
		fire := animations.NewFireEffect(width, height, palette)
		return &AnimationWrapper{
			render: fire.Render,
			update: fire.Update,
		}

	case "matrix":
		palette := animations.GetMatrixPalette(themeName)
		matrix := animations.NewMatrixEffect(width, height, palette)
		return &AnimationWrapper{
			render: matrix.Render,
			update: matrix.Update,
		}

	case "matrix-art":
		palette := animations.GetMatrixPalette(themeName)
		text := m.loadTextFile(fileName)
		matrixArt := animations.NewMatrixArtEffect(width, height, palette, text)
		return &AnimationWrapper{
			render: matrixArt.Render,
			update: matrixArt.Update,
		}

	case "rain":
		palette := animations.GetRainPalette(themeName)
		rain := animations.NewRainEffect(width, height, palette)
		return &AnimationWrapper{
			render: rain.Render,
			update: rain.Update,
		}

	case "rain-art":
		palette := animations.GetRainPalette(themeName)
		text := m.loadTextFile(fileName)
		rainArt := animations.NewRainArtEffect(width, height, palette, text)
		return &AnimationWrapper{
			render: rainArt.Render,
			update: rainArt.Update,
		}

	case "fireworks":
		palette := animations.GetFireworksPalette(themeName)
		fireworks := animations.NewFireworksEffect(width, height, palette)
		return &AnimationWrapper{
			render: fireworks.Render,
			update: fireworks.Update,
		}

	case "pour":
		text := m.loadTextFile(fileName)
		config := animations.PourConfig{
			Width:                  width,
			Height:                 height,
			Text:                   text,
			PourDirection:          "down",
			PourSpeed:              3,
			MovementSpeed:          0.2,
			Gap:                    1,
			StartingColor:          "#ffffff",
			FinalGradientStops:     getGradientStops(themeName),
			FinalGradientSteps:     12,
			FinalGradientFrames:    5,
			FinalGradientDirection: "horizontal",
		}
		pour := animations.NewPourEffect(config)
		return &AnimationWrapper{
			render: pour.Render,
			update: pour.Update,
		}

	case "print":
		text := m.loadTextFile(fileName)
		config := animations.PrintConfig{
			Width:           width,
			Height:          height,
			Text:            text,
			CharDelay:       30 * time.Millisecond,
			PrintSpeed:      2,
			PrintHeadSymbol: "█",
			TrailSymbols:    []string{"░", "▒", "▓"},
			GradientStops:   getGradientStops(themeName),
		}
		print := animations.NewPrintEffect(config)
		return &AnimationWrapper{
			render: print.Render,
			update: print.Update,
		}

	// TODO: Config-based animations below need full config structs
	// For now, return nil and they will fall back to external launch
	case "beams", "beam-text", "ring-text", "blackhole-text", "aquarium":
		return nil

	default:
		// Unsupported animation type - return nil
		return nil
	}
}

// loadTextFile loads a text file for text-based animations
func (m *Model) loadTextFile(filename string) string {
	if filename == "" || filename == "(n/a)" {
		return "SYSC" // fallback
	}

	filePath := getAssetPath(filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "SYSC" // fallback on error
	}

	return string(data)
}

// getGradientStops returns gradient color stops for the given theme
func getGradientStops(theme string) []string {
	switch theme {
	case "dracula":
		return []string{"#ff79c6", "#bd93f9", "#ffffff"}
	case "gruvbox":
		return []string{"#fe8019", "#fabd2f", "#ffffff"}
	case "nord":
		return []string{"#88c0d0", "#81a1c1", "#ffffff"}
	case "tokyo-night":
		return []string{"#9ece6a", "#e0af68", "#ffffff"}
	case "catppuccin":
		return []string{"#cba6f7", "#f5c2e7", "#ffffff"}
	case "material":
		return []string{"#03dac6", "#bb86fc", "#ffffff"}
	case "solarized":
		return []string{"#268bd2", "#2aa198", "#ffffff"}
	case "monochrome":
		return []string{"#808080", "#c0c0c0", "#ffffff"}
	case "transishardjob":
		return []string{"#55cdfc", "#f7a8b8", "#ffffff"}
	case "rama":
		return []string{"#ef233c", "#d90429", "#edf2f4"}
	case "eldritch":
		return []string{"#37f499", "#04d1f9", "#ebfafa"}
	case "dark":
		return []string{"#ffffff", "#cccccc", "#ffffff"}
	default:
		return []string{"#8A008A", "#00D1FF", "#FFFFFF"}
	}
}

// getBeamColors returns beam colors for the given theme
func getBeamColors(theme string) []string {
	switch theme {
	case "dracula":
		return []string{"#ff79c6", "#bd93f9", "#8be9fd", "#50fa7b", "#ffb86c"}
	case "gruvbox":
		return []string{"#fb4934", "#fe8019", "#fabd2f", "#b8bb26", "#83a598"}
	case "nord":
		return []string{"#bf616a", "#d08770", "#ebcb8b", "#a3be8c", "#88c0d0"}
	case "tokyo-night":
		return []string{"#f7768e", "#ff9e64", "#e0af68", "#9ece6a", "#73daca"}
	case "catppuccin":
		return []string{"#f38ba8", "#fab387", "#f9e2af", "#a6e3a1", "#89dceb"}
	case "material":
		return []string{"#f07178", "#ff9cac", "#03dac6", "#bb86fc", "#ff6e40"}
	case "solarized":
		return []string{"#dc322f", "#cb4b16", "#b58900", "#859900", "#268bd2"}
	case "monochrome":
		return []string{"#ffffff", "#d0d0d0", "#a0a0a0", "#808080", "#606060"}
	case "transishardjob":
		return []string{"#55cdfc", "#f7a8b8", "#ffffff", "#f7a8b8", "#55cdfc"}
	case "rama":
		return []string{"#ef233c", "#d90429", "#8d99ae", "#2b2d42", "#edf2f4"}
	case "eldritch":
		return []string{"#37f499", "#04d1f9", "#f7c67f", "#f16c75", "#ebfafa"}
	case "dark":
		return []string{"#ffffff", "#cccccc", "#999999", "#666666", "#444444"}
	default:
		return []string{"#FF0080", "#8A008A", "#00D1FF", "#00FF00", "#FFFF00"}
	}
}

// getAquariumColors returns aquarium colors for the given theme
func getAquariumColors(theme string) []string {
	switch theme {
	case "dracula":
		return []string{"#ff79c6", "#bd93f9", "#8be9fd"}
	case "gruvbox":
		return []string{"#fe8019", "#b8bb26", "#83a598"}
	case "nord":
		return []string{"#88c0d0", "#81a1c1", "#5e81ac"}
	case "tokyo-night":
		return []string{"#73daca", "#7aa2f7", "#9ece6a"}
	case "catppuccin":
		return []string{"#89dceb", "#89b4fa", "#cba6f7"}
	case "material":
		return []string{"#03dac6", "#bb86fc", "#018786"}
	case "solarized":
		return []string{"#268bd2", "#2aa198", "#859900"}
	case "monochrome":
		return []string{"#ffffff", "#c0c0c0", "#808080"}
	case "transishardjob":
		return []string{"#55cdfc", "#f7a8b8", "#ffffff"}
	case "rama":
		return []string{"#8d99ae", "#edf2f4", "#ef233c"}
	case "eldritch":
		return []string{"#04d1f9", "#37f499", "#a48cf4"}
	case "dark":
		return []string{"#ffffff", "#cccccc", "#999999"}
	default:
		return []string{"#00D1FF", "#8A008A", "#00FF00"}
	}
}
