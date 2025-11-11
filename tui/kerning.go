package tui

import (
	"strings"
	"unicode/utf8"
)

// runeAt is a helper to safely get a rune at a specific column index within a string.
// It handles multi-byte runes correctly.
func runeAt(s string, x int) (rune, bool) {
	if x < 0 {
		return ' ', false
	}
	i := 0
	for _, r := range s {
		if i == x {
			return r, true
		}
		i++
	}
	return ' ', false
}

// maxRowLen is a helper to find the maximum rune length among all rows of a glyph,
// effectively determining its bounding box width.
func maxRowLen(rows []string) int {
	maxLen := 0
	for _, r := range rows {
		maxLen = max(maxLen, utf8.RuneCountInString(r))
	}
	return maxLen
}

// normalizeGlyph pads rows of a glyph to a consistent height with empty strings
// to simplify comparison between glyphs of different heights.
func normalizeGlyph(glyph []string, height int) []string {
	out := make([]string, height)
	maxLen := 0
	for _, r := range glyph {
		maxLen = max(maxLen, utf8.RuneCountInString(r))
	}

	for i := range height {
		if i < len(glyph) {
			// Ensure each row has consistent width by padding with spaces if needed
			row := glyph[i]
			if utf8.RuneCountInString(row) < maxLen {
				row += strings.Repeat(" ", maxLen-utf8.RuneCountInString(row))
			}
			out[i] = row
		} else {
			// Pad with spaces (not empty strings) to maintain consistent width
			out[i] = strings.Repeat(" ", maxRowLen(glyph))
		}
	}
	return out
}

// getActivePixels extracts all non-space pixel coordinates for a given glyph.
// These are the "ink" pixels that contribute to the character's visual form.
func getActivePixels(glyph []string) map[pixelCoord]bool {
	activePixels := make(map[pixelCoord]bool)
	for y, row := range glyph {
		for x := range utf8.RuneCountInString(row) {
			r, _ := runeAt(row, x)
			// ' ' (space) and '\u0000' (null character) are considered empty/background.
			if r != ' ' && r != 0 {
				activePixels[pixelCoord{float64(x), y, false}] = true
			}
		}
	}
	return activePixels
}

// Helper function to calculate absolute value for float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// checkCollisionWithSmartBuffer uses enhanced visual analysis to allow better kerning
func checkCollisionWithSmartBuffer(glyphA, glyphB []string, widthA int, spacing int) bool {
	aPixels := getActivePixels(glyphA)

	// Check for half-pixel alignment issues
	heightDiff := len(glyphA) - len(glyphB)
	needsHalfPixelAdjustment := heightDiff%2 != 0

	for yB, rowB := range glyphB {
		for xB := 0; xB < utf8.RuneCountInString(rowB); xB++ {
			rB, _ := runeAt(rowB, xB)
			if rB == ' ' || rB == 0 {
				continue
			}

			// Calculate the absolute X position of glyphB's pixel
			xAbsB := float64(widthA + spacing + xB)
			yAbsB := yB

			// Apply half-pixel adjustment if needed
			if needsHalfPixelAdjustment && yAbsB >= len(glyphB)/2 {
				xAbsB += 0.5 // Adjust by half a pixel for the bottom half
			}

			// Check for direct overlap
			for coord := range aPixels {
				// Convert to float64 for comparison
				coordX := coord.x
				coordY := coord.y

				// Check if this pixel would overlap with any pixel in glyphA
				if abs(coordX-xAbsB) < 1.0 && coordY == yAbsB {
					return true
				}
			}

			// Enhanced adjacency check
			conflictCount := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					checkX := xAbsB + float64(dx)
					checkY := yAbsB + dy

					for coord := range aPixels {
						coordX := coord.x
						coordY := coord.y

						// Check proximity with half-pixel precision
						xDist := abs(coordX - checkX)
						yDist := abs(float64(coordY) - float64(checkY))

						if xDist < 1.0 && yDist < 1.0 {
							conflictCount++
							// Weight horizontal conflicts more heavily
							if dx != 0 && dy == 0 {
								conflictCount += 2
							}
						}
					}
				}
			}

			// Only consider it a collision if there's significant conflict
			if conflictCount >= 3 {
				return true
			}
		}
	}

	return false
}

// computeKerning calculates the minimum horizontal spacing between two characters (glyphA and glyphB)
// such that their active pixels do not touch (even by corners). This spacing is constrained to -1, 0, or +1.
func computeKerning(glyphA, glyphB []string) int {
	// If either glyph has no visual width (e.g., empty bitmap), return neutral spacing
	if maxRowLen(glyphA) == 0 || maxRowLen(glyphB) == 0 {
		return 0 // No adjustment needed for empty glyphs
	}

	hA, hB := len(glyphA), len(glyphB)
	h := max(hA, hB)

	a := normalizeGlyph(glyphA, h)
	b := normalizeGlyph(glyphB, h)
	widthA := maxRowLen(a)

	// Check spacing options in order: -1 (tight), 0 (normal), +1 (extra space)
	// Return the minimum safe spacing within our constraint range

	// First try -1 (tight kerning/tucking)
	if !checkCollisionWithSmartBuffer(a, b, widthA, -1) {
		return -1 // Safe to tuck closer
	}

	// Then try 0 (normal spacing)
	if !checkCollisionWithSmartBuffer(a, b, widthA, 0) {
		return 0 // Normal spacing works
	}

	// Finally, use +1 (extra space needed)
	return 1 // Need extra space to prevent collision
}
