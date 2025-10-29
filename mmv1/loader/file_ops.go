package loader

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func GlobWithBase(baseDir, pattern string) ([]string, error) {
	if baseDir == "" {
		return filepath.Glob(pattern)
	}

	// Ensure clean concatenation
	if !strings.HasSuffix(baseDir, "/") {
		baseDir = baseDir + "/"
	}

	// Simple concatenation preserves patterns
	fullPattern := baseDir + pattern
	return filepath.Glob(fullPattern)
}

func Exists(paths ...string) bool {
	fullPath := filepath.Join(paths...)
	_, err := os.Stat(fullPath)
	exists := !errors.Is(err, os.ErrNotExist)
	return exists
}
