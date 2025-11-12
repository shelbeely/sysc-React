package tui

import (
	"os"

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

	// TODO: These animations use Config structs and need refactoring
	// case "pour", "print", "beams", "beam-text", "ring-text", "blackhole-text", "aquarium":
	//   return nil  // Not yet supported in viewport mode

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
