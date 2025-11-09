package tui

import (
	"os"
	"path/filepath"
	"strings"
)

// discoverAssetFiles finds all .txt files in the assets directory
func discoverAssetFiles() []string {
	var files []string

	// Try multiple possible asset paths
	assetPaths := []string{
		"assets",
		"./assets",
		"../assets",
		filepath.Join(os.Getenv("HOME"), "sysc-Go", "assets"),
		"/usr/local/share/sysc-go/assets",
	}

	for _, assetPath := range assetPaths {
		entries, err := os.ReadDir(assetPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasSuffix(strings.ToLower(name), ".txt") {
				files = append(files, name)
			}
		}

		// If we found files, return them
		if len(files) > 0 {
			return files
		}
	}

	return files
}

// getAssetPath returns the full path to an asset file
func getAssetPath(filename string) string {
	assetPaths := []string{
		filepath.Join("assets", filename),
		filepath.Join("./assets", filename),
		filepath.Join("../assets", filename),
		filepath.Join(os.Getenv("HOME"), "sysc-Go", "assets", filename),
		filepath.Join("/usr/local/share/sysc-go/assets", filename),
	}

	for _, path := range assetPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filename // fallback
}
