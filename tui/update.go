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
		// Update textarea size if in editor mode
		if m.editorMode {
			m.textarea.SetWidth(m.width - 10)
			m.textarea.SetHeight(m.height - 10)
		}
		return m, nil

	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle editor mode separately
	if m.editorMode {
		return m.handleEditorKeyPress(msg)
	}

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

	// Check if "Custom text" is selected - enter editor mode instead
	if fileName == "Custom text" {
		m.editorMode = true
		m.textarea.Focus()
		return m, nil
	}

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

// handleEditorKeyPress handles keyboard input in editor mode
func (m Model) handleEditorKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle export prompt separately
	if m.showExportPrompt {
		switch msg.String() {
		case "esc":
			// Cancel export prompt
			m.showExportPrompt = false
			return m, nil

		case "up", "k":
			// Navigate up in export options
			if m.exportTarget > 0 {
				m.exportTarget--
			}
			return m, nil

		case "down", "j":
			// Navigate down in export options
			if m.exportTarget < 1 { // 0=syscgo, 1=sysc-walls
				m.exportTarget++
			}
			return m, nil

		case "enter":
			// Confirm export target
			if m.exportTarget == 1 {
				// sysc-walls is not implemented yet
				// TODO: Show a message that this feature is WIP
				m.showExportPrompt = false
				return m, nil
			}
			// syscgo export - show save prompt
			m.showExportPrompt = false
			m.showSavePrompt = true
			m.filenameInput.Focus()
			return m, nil
		}
		return m, nil
	}

	// Handle save prompt separately
	if m.showSavePrompt {
		switch msg.String() {
		case "esc":
			// Cancel save prompt
			m.showSavePrompt = false
			m.filenameInput.SetValue("")
			m.filenameInput.Blur()
			return m, nil

		case "enter":
			// Save file
			return m.saveFile()

		default:
			// Update filename input
			m.filenameInput, cmd = m.filenameInput.Update(msg)
			return m, cmd
		}
	}

	// Handle editor input
	switch msg.String() {
	case "esc":
		// Exit editor mode and return to main UI
		m.editorMode = false
		m.textarea.Blur()
		m.textarea.Reset()
		// Refresh file list to include any newly saved files
		files := discoverAssetFiles()
		if len(files) == 0 {
			files = []string{"SYSC.txt"}
		}
		files = append([]string{"Custom text"}, files...)
		m.files = files
		return m, nil

	case "ctrl+s":
		// Show export/save prompt to choose target
		m.showExportPrompt = true
		m.exportTarget = 0 // Default to syscgo
		return m, nil

	default:
		// Update textarea
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}
}

// saveFile saves the text area content to assets folder
func (m Model) saveFile() (Model, tea.Cmd) {
	// Clear previous error
	m.saveError = ""

	filename := m.filenameInput.Value()
	if filename == "" {
		m.saveError = "Filename cannot be empty"
		return m, nil
	}

	// Add .txt extension if not present
	if len(filename) < 4 || filename[len(filename)-4:] != ".txt" {
		filename += ".txt"
	}

	// Validate content is not empty
	if len(m.textarea.Value()) == 0 {
		m.saveError = "Content cannot be empty"
		return m, nil
	}

	// Save to assets folder
	err := saveToAssets(filename, m.textarea.Value())
	if err != nil {
		m.saveError = err.Error()
		return m, nil
	}

	// Exit editor mode and return to main UI
	m.editorMode = false
	m.showSavePrompt = false
	m.saveError = ""
	m.textarea.Blur()
	m.textarea.Reset()
	m.filenameInput.SetValue("")
	m.filenameInput.Blur()

	// Refresh file list to include the newly saved file
	files := discoverAssetFiles()
	if len(files) == 0 {
		files = []string{"SYSC.txt"}
	}
	files = append([]string{"Custom text"}, files...)
	m.files = files

	// Set the selected file to the newly created file
	for i, file := range m.files {
		if file == filename {
			m.selectedFile = i
			break
		}
	}

	return m, nil
}
