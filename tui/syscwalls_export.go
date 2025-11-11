package tui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportToSyscWalls exports ASCII art to sysc-walls with config update
func ExportToSyscWalls(filename, content string) error {
	// Create walls directory
	wallsDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "syscgo", "walls")
	if err := os.MkdirAll(wallsDir, 0755); err != nil {
		return fmt.Errorf("failed to create walls directory: %w", err)
	}

	// Save ASCII art file
	artPath := filepath.Join(wallsDir, filename)
	if err := os.WriteFile(artPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to save ASCII art: %w", err)
	}

	// Update daemon.conf
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "sysc-walls", "daemon.conf")
	if err := updateSyscWallsConfig(configPath, artPath); err != nil {
		// Non-fatal - file saved successfully
		return fmt.Errorf("ASCII art saved to %s, but failed to update config: %w", artPath, err)
	}

	return nil
}

// updateSyscWallsConfig updates or creates the sysc-walls daemon config
func updateSyscWallsConfig(configPath, artPath string) error {
	// Create config directory if needed
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config or create default
	config := make(map[string]map[string]string)
	config["idle"] = map[string]string{"timeout": "300"}
	config["daemon"] = map[string]string{"monitors": "all"}
	config["terminal"] = map[string]string{"fullscreen": "true", "opacity": "0.95"}
	config["animation"] = map[string]string{
		"type":     "beam-text",
		"theme":    "dracula",
		"file":     artPath,
		"duration": "infinite",
	}

	// If config exists, read and merge
	if data, err := os.ReadFile(configPath); err == nil {
		parseINI(string(data), config)
		// Update animation section with new file
		if config["animation"] == nil {
			config["animation"] = make(map[string]string)
		}
		config["animation"]["file"] = artPath
		config["animation"]["type"] = "beam-text"
		config["animation"]["theme"] = "dracula"
		config["animation"]["duration"] = "infinite"
	}

	// Write config
	if err := writeINI(configPath, config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// parseINI parses a simple INI format config
func parseINI(content string, config map[string]map[string]string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			if config[currentSection] == nil {
				config[currentSection] = make(map[string]string)
			}
			continue
		}

		// Key-value pair
		if currentSection != "" && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[currentSection][key] = value
		}
	}
}

// writeINI writes config to INI format
func writeINI(path string, config map[string]map[string]string) error {
	var content strings.Builder

	// Write sections in a specific order
	sectionOrder := []string{"idle", "daemon", "animation", "terminal"}

	for _, section := range sectionOrder {
		if values, ok := config[section]; ok {
			content.WriteString(fmt.Sprintf("[%s]\n", section))
			for key, value := range values {
				content.WriteString(fmt.Sprintf("%s = %s\n", key, value))
			}
			content.WriteString("\n")
			delete(config, section) // Mark as written
		}
	}

	// Write any remaining sections
	for section, values := range config {
		content.WriteString(fmt.Sprintf("[%s]\n", section))
		for key, value := range values {
			content.WriteString(fmt.Sprintf("%s = %s\n", key, value))
		}
		content.WriteString("\n")
	}

	return os.WriteFile(path, []byte(content.String()), 0644)
}

// ExportBitArt handles export target selection and saves accordingly
func ExportBitArt(filename string, content []string, target int) error {
	// Strip ANSI codes
	plainContent := ""
	for _, line := range content {
		plainContent += stripANSI(line) + "\n"
	}

	// Add .txt extension if not present
	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}

	switch target {
	case 0: // syscgo
		return saveToAssets(filename, plainContent)

	case 1: // sysc-walls
		return ExportToSyscWalls(filename, plainContent)

	default:
		return fmt.Errorf("unknown export target: %d", target)
	}
}
