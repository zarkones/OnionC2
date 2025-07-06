package core

import (
	"path/filepath"
	"strings"
)

func Deduplicate[T string](slice []T) []T {
	ids := make(map[T]bool)
	l := []T{}
	for _, item := range slice {
		if _, v := ids[item]; !v {
			ids[item] = true
			l = append(l, item)
		}
	}
	return l
}

// GetFileName returns file's name from a path.
func GetFileName(path string) string {
	// Replace all backslashes with forward slashes.
	normalizedPath := strings.ReplaceAll(path, "\\", "/")
	// Clean the path to handle relative components.
	cleanPath := filepath.Clean(normalizedPath)
	// Extract the base name.
	fileName := filepath.Base(cleanPath)
	// If the path still contains a drive letter (e.g., "C:"), take the last segment after the last "/".
	if idx := strings.Index(fileName, ":"); idx != -1 {
		fileName = filepath.Base(strings.TrimPrefix(cleanPath, fileName[:idx+1]))
	}
	return fileName
}
