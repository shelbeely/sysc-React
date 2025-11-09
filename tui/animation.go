package tui

import (
	"fmt"
	"strings"
)

// CreateAnimationPreview creates a static preview of the animation config
func CreateAnimationPreview(animName, theme, file string, width, height int) string {
	var preview strings.Builder

	// Center the preview text
	centerY := height / 2
	emptyLines := centerY - 6

	for i := 0; i < emptyLines; i++ {
		preview.WriteString("\n")
	}

	// Animation info
	preview.WriteString(centerText("╔════════════════════════════════╗", width))
	preview.WriteString("\n")
	preview.WriteString(centerText("║    ANIMATION READY TO START    ║", width))
	preview.WriteString("\n")
	preview.WriteString(centerText("╚════════════════════════════════╝", width))
	preview.WriteString("\n\n")

	preview.WriteString(centerText(fmt.Sprintf("Animation: %s", animName), width))
	preview.WriteString("\n")
	preview.WriteString(centerText(fmt.Sprintf("Theme: %s", theme), width))
	preview.WriteString("\n")
	preview.WriteString(centerText(fmt.Sprintf("File: %s", file), width))
	preview.WriteString("\n\n")

	preview.WriteString(centerText("Press ENTER to launch animation", width))
	preview.WriteString("\n")
	preview.WriteString(centerText("(opens in new view)", width))

	return preview.String()
}

// centerText centers text within the given width
func centerText(text string, width int) string {
	// Remove ANSI codes for length calculation
	plainText := stripANSI(text)
	textLen := len(plainText)

	if textLen >= width {
		return text
	}

	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}

// stripANSI removes ANSI escape codes from text
func stripANSI(text string) string {
	// Simple ANSI stripper - good enough for our purposes
	inEscape := false
	var result strings.Builder

	for _, ch := range text {
		if ch == '\033' {
			inEscape = true
		} else if inEscape {
			if ch == 'm' {
				inEscape = false
			}
		} else {
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// LaunchAnimation launches the actual animation in the CLI
func LaunchAnimation(animName, theme, file, duration string) error {
	// For the MVP, we just return - the animation parameters are set
	// The user will see the TUI quit and can manually run syscgo if needed
	// TODO: Actually exec syscgo with the selected parameters using os/exec
	// This would require:
	// 1. Finding the syscgo binary
	// 2. Building the command with proper flags
	// 3. Executing it with os.Exec or similar
	return nil
}
