package tui

import (
	"strings"
)

// VerticalAlignment represents the vertical alignment of a character or line.
type VerticalAlignment int

const (
	// AlignTop aligns the top of the character/line to the top of the canvas.
	AlignTop VerticalAlignment = iota
	// AlignMiddle aligns the middle of the character/line to the middle of the canvas.
	AlignMiddle
	// AlignBottom aligns the bottom of the character/line to the bottom of the canvas.
	AlignBottom
)

// HorizontalAlignment represents the horizontal alignment of a character or line.
type HorizontalAlignment int

const (
	// AlignLeft aligns the left of the character/line to the left of the canvas.
	AlignLeft HorizontalAlignment = iota
	// AlignCenter aligns the center of the character/line to the center of the canvas.
	AlignCenter
	// AlignRight aligns the right of the character/line to the right of the canvas.
	AlignRight
)

// ApplyVerticalAlignment applies vertical alignment to a set of lines.
func ApplyVerticalAlignment(lines []string, canvasHeight int, alignment VerticalAlignment) []string {
	if len(lines) >= canvasHeight {
		return lines
	}

	padding := canvasHeight - len(lines)
	var topPadding int

	switch alignment {
	case AlignTop:
		// no padding needed
	case AlignMiddle:
		topPadding = padding / 2
	case AlignBottom:
		topPadding = padding
	}

	result := make([]string, canvasHeight)
	// Top padding is already zero-initialized (empty strings)
	copy(result[topPadding:], lines)
	// Bottom padding is already zero-initialized (empty strings)

	return result
}

// ApplyHorizontalAlignment applies horizontal alignment to a single line of text.
func ApplyHorizontalAlignment(line string, canvasWidth int, alignment HorizontalAlignment) string {
	lineWidth := len(line)
	if lineWidth >= canvasWidth {
		return line
	}

	padding := canvasWidth - lineWidth
	var leftPadding, rightPadding int

	switch alignment {
	case AlignLeft:
		rightPadding = padding
	case AlignCenter:
		leftPadding = padding / 2
		rightPadding = padding - leftPadding
	case AlignRight:
		leftPadding = padding
	}

	return strings.Repeat(" ", leftPadding) + line + strings.Repeat(" ", rightPadding)
}

// analyzeDescenderProperties analyzes the font data to determine descender properties for each character
func analyzeDescenderProperties(fontData FontData, scaleFactor float64) map[string]DescenderInfo {
	descenderMap := make(map[string]DescenderInfo)

	// First pass: find the common baseline by analyzing lowercase letters without descenders
	baselineRow := findCommonBaseline(fontData, scaleFactor)

	// Second pass: analyze each character relative to the common baseline
	for charStr, bitmapLines := range fontData.Characters {
		scaledLines := scaleCharacter(bitmapLines, scaleFactor)
		info := analyzeCharacterDescenders(scaledLines, baselineRow)
		descenderMap[charStr] = info
	}

	return descenderMap
}

// findCommonBaseline determines the baseline position by analyzing lowercase letters without descenders
func findCommonBaseline(fontData FontData, scaleFactor float64) int {
	// Sample lowercase letters that typically don't have descenders
	sampleChars := []string{"a", "e", "o", "x", "n", "m", "s", "c"}

	var baselinePositions []int

	for _, charStr := range sampleChars {
		if bitmapLines, ok := fontData.Characters[charStr]; ok {
			scaledLines := scaleCharacter(bitmapLines, scaleFactor)

			// Find the last row with content (this is the baseline for non-descender chars)
			lastContentRow := -1
			for row := len(scaledLines) - 1; row >= 0; row-- {
				hasContent := false
				for _, char := range scaledLines[row] {
					if char != ' ' {
						hasContent = true
						break
					}
				}
				if hasContent {
					lastContentRow = row
					break
				}
			}

			if lastContentRow >= 0 {
				baselinePositions = append(baselinePositions, lastContentRow)
			}
		}
	}

	// Use the most common baseline position (or average if they vary)
	if len(baselinePositions) == 0 {
		return 5 // Fallback default
	}

	// Calculate average baseline
	sum := 0
	for _, pos := range baselinePositions {
		sum += pos
	}
	return sum / len(baselinePositions)
}

// analyzeCharacterDescenders analyzes a single character's bitmap to determine its descender properties
func analyzeCharacterDescenders(bitmapLines []string, commonBaseline int) DescenderInfo {
	if len(bitmapLines) == 0 {
		return DescenderInfo{HasDescender: false, BaselineHeight: 0, DescenderHeight: 0, TotalHeight: 0, VerticalOffset: 0}
	}

	height := len(bitmapLines)

	// Find the first and last rows with content
	firstContentRow := -1
	lastContentRow := -1

	for row := range height {
		hasContent := false
		for _, char := range bitmapLines[row] {
			if char != ' ' {
				hasContent = true
				break
			}
		}
		if hasContent {
			if firstContentRow == -1 {
				firstContentRow = row
			}
			lastContentRow = row
		}
	}

	if lastContentRow == -1 {
		// Empty character (like space)
		return DescenderInfo{HasDescender: false, BaselineHeight: height, DescenderHeight: 0, TotalHeight: height, VerticalOffset: 0}
	}

	// A character has a descender if its last content row extends beyond the common baseline
	hasDescender := lastContentRow > commonBaseline
	descenderHeight := 0
	if hasDescender {
		descenderHeight = lastContentRow - commonBaseline
	}

	return DescenderInfo{
		HasDescender:    hasDescender,
		BaselineHeight:  commonBaseline + 1,
		DescenderHeight: descenderHeight,
		TotalHeight:     height,
		VerticalOffset:  0,
	}
}

// adjustCharacterForDescenders adjusts a character's bitmap for proper descender alignment
func adjustCharacterForDescenders(bitmapLines []string, info DescenderInfo, maxCharHeight int) []string {
	// All characters should be aligned to the common baseline
	// Characters with descenders that are shorter than maxCharHeight need padding at the bottom
	// Characters without descenders also get padded at the bottom

	if len(bitmapLines) >= maxCharHeight {
		return bitmapLines
	}

	// Calculate padding needed at the bottom to reach max height
	paddingNeeded := maxCharHeight - len(bitmapLines)
	if paddingNeeded <= 0 {
		return bitmapLines
	}

	// Add padding at the bottom for all characters
	result := make([]string, maxCharHeight)
	copy(result, bitmapLines)

	// Fill remaining rows with empty strings (no spaces, just empty)
	// This ensures the width calculation isn't affected
	for i := len(bitmapLines); i < maxCharHeight; i++ {
		result[i] = ""
	}

	return result
}
