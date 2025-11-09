package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Canvas takes up most of the screen, leave room for selectors and help
		m.canvasHeight = m.height - 12
		return m, nil

	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "esc":
		return m, tea.Quit

	case "up", "k":
		return m.navigateUp(), nil

	case "down", "j":
		return m.navigateDown(), nil

	case "left", "h":
		return m.navigateLeft(), nil

	case "right", "l":
		return m.navigateRight(), nil

	case "enter":
		return m.startAnimation()
	}

	return m, nil
}

// navigateUp moves the selection up within the current selector
func (m Model) navigateUp() Model {
	switch m.focusedSelector {
	case 0: // Animation selector
		if m.selectedAnimation > 0 {
			m.selectedAnimation--
		}
	case 1: // Theme selector
		if m.selectedTheme > 0 {
			m.selectedTheme--
		}
	case 2: // File selector
		if m.selectedFile > 0 {
			m.selectedFile--
		}
	case 3: // Duration selector
		if m.selectedDuration > 0 {
			m.selectedDuration--
		}
	}
	return m
}

// navigateDown moves the selection down within the current selector
func (m Model) navigateDown() Model {
	switch m.focusedSelector {
	case 0: // Animation selector
		if m.selectedAnimation < len(m.animations)-1 {
			m.selectedAnimation++
		}
	case 1: // Theme selector
		if m.selectedTheme < len(m.themes)-1 {
			m.selectedTheme++
		}
	case 2: // File selector
		if m.selectedFile < len(m.files)-1 {
			m.selectedFile++
		}
	case 3: // Duration selector
		if m.selectedDuration < len(m.durations)-1 {
			m.selectedDuration++
		}
	}
	return m
}

// navigateLeft moves focus to the previous selector
func (m Model) navigateLeft() Model {
	if m.focusedSelector > 0 {
		m.focusedSelector--
	}
	return m
}

// navigateRight moves focus to the next selector
func (m Model) navigateRight() Model {
	if m.focusedSelector < 3 {
		m.focusedSelector++
	}
	return m
}

// startAnimation begins the animation with current settings
func (m Model) startAnimation() (Model, tea.Cmd) {
	// Get selected values
	animName := m.animations[m.selectedAnimation]
	theme := m.themes[m.selectedTheme]
	fileName := m.files[m.selectedFile]
	duration := m.durations[m.selectedDuration]

	// Launch animation immediately and quit TUI
	// This returns user to normal terminal where animation runs
	return m, tea.Sequence(
		tea.Quit,
		func() tea.Msg {
			// Launch animation after TUI quits
			LaunchAnimation(animName, theme, fileName, duration)
			return nil
		},
	)
}
