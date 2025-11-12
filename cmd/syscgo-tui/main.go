package main

import (
	"fmt"
	"os"

	"github.com/Nomadcxx/sysc-Go/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the TUI model
	m := tui.NewModel()

	// Create the program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
