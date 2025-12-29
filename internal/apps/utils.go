package apps

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateID creates a simple hash for identifying apps
func generateID(input string) string {
	hash := 0
	for _, c := range input {
		hash = hash*31 + int(c)
	}
	return fmt.Sprintf("app_%d", uint32(hash))
}

// findExecutables finds all .exe files in a directory (max 2 levels deep)
func findExecutables(installPath string) []string {
	var exes []string

	filepath.Walk(installPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(installPath, path)
		if strings.Count(rel, string(filepath.Separator)) > 2 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".exe") {
			exes = append(exes, path)
		}
		return nil
	})

	return exes
}

// extractIconBase64 attempts to find and encode an icon for the app
func extractIconBase64(installPath string) string {
	if installPath == "" {
		return ""
	}

	patterns := []string{"*Logo*.png", "*logo*.png", "Assets/*Logo*.png", "Assets/*logo*.png"}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(installPath, pattern))
		if len(matches) > 0 {
			if data, err := os.ReadFile(matches[0]); err == nil {
				return base64.StdEncoding.EncodeToString(data)
			}
		}
	}
	return ""
}

// cleanStoreName makes Store app names more readable
func cleanStoreName(name string) string {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 && len(parts[0]) > 8 {
		return parts[1]
	}
	return name
}

// cleanPublisher removes "CN=" prefix from publisher
func cleanPublisher(publisher string) string {
	publisher = strings.TrimPrefix(publisher, "CN=")
	if idx := strings.Index(publisher, ","); idx > 0 {
		publisher = publisher[:idx]
	}
	return publisher
}

// isSystemApp filters out Windows system components
func isSystemApp(name, publisher, installPath string) bool {
	lowerName := strings.ToLower(name)
	lowerPublisher := strings.ToLower(publisher)
	lowerPath := strings.ToLower(installPath)

	if strings.Contains(lowerPublisher, "microsoft") {
		allowed := []string{"office", "visual studio", "vscode", "edge", "teams", "onedrive"}
		for _, a := range allowed {
			if strings.Contains(lowerName, a) {
				return false
			}
		}
		system := []string{"update", "redistributable", "runtime", ".net", "sdk", "tool"}
		for _, s := range system {
			if strings.Contains(lowerName, s) {
				return true
			}
		}
	}

	if strings.Contains(lowerPath, "\\windows\\") {
		return true
	}
	return false
}

// isSystemStoreApp filters out Windows system Store apps
func isSystemStoreApp(name, publisher string) bool {
	lowerName := strings.ToLower(name)
	prefixes := []string{
		"microsoft.net", "microsoft.vclibs", "microsoft.ui",
		"microsoft.windows", "microsoft.services", "microsoft.advertising",
		"microsoft.directx", "microsoft.desktop",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lowerName, prefix) {
			return true
		}
	}
	return false
}
