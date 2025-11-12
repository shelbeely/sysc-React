package main

import (
	"fmt"
	"strings"
	"github.com/Nomadcxx/sysc-Go/animations"
)

func main() {
	width := 80
	height := 24

	// Test ring-text
	ringColors := []string{"#bd93f9", "#ff79c6", "#f1fa8c", "#8be9fd", "#50fa7b", "#ffb86c"}
	config := animations.RingTextConfig{
		Width: width,
		Height: height,
		Text: "TEST",
		RingColors: ringColors,
		RingGap: 0.1,
		SpinSpeedRange: [2]float64{0.025, 0.075},
		SpinDuration: 200,
		DisperseDuration: 200,
		SpinDisperseCycles: 3,
		TransitionFrames: 60,
		StaticFrames: 30,
		FinalGradientStops: []string{"#6272a4", "#bd93f9", "#f8f8f2"},
		FinalGradientSteps: 12,
		StaticGradientStops: ringColors,
		StaticGradientDir: animations.GradientHorizontal,
	}

	ringText := animations.NewRingTextEffect(config)

	// Update a few times
	for i := 0; i < 50; i++ {
		ringText.Update()
	}

	output := ringText.Render()

	// Count lines
	lines := strings.Split(output, "\n")

	fmt.Printf("Terminal dimensions: %dx%d\n", width, height)
	fmt.Printf("Render() returned %d lines\n", len(lines))
	fmt.Printf("Expected: %d lines\n", height)

	if len(lines) != height {
		fmt.Printf("\n⚠️ MISMATCH! Animation returns %d lines but terminal is %d lines high\n", len(lines), height)
		fmt.Printf("This causes gaps when printing to CLI!\n")
	} else {
		fmt.Printf("\n✓ Line count matches terminal height\n")
	}

	// Check for empty lines
	emptyCount := 0
	for i, line := range lines {
		plain := stripANSI(line)
		if strings.TrimSpace(plain) == "" {
			emptyCount++
			if i < 5 || i >= len(lines)-5 {
				fmt.Printf("Line %d is empty\n", i)
			}
		}
	}
	fmt.Printf("\nEmpty lines: %d (%.1f%%)\n", emptyCount, float64(emptyCount)/float64(len(lines))*100)
}

func stripANSI(s string) string {
	inAnsi := false
	var result strings.Builder
	for _, r := range s {
		if r == '\033' {
			inAnsi = true
			continue
		}
		if inAnsi {
			if r == 'm' {
				inAnsi = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}
