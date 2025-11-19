package loader

import (
	"errors"
	"os"
	"path/filepath"
)

func Exists(paths ...string) bool {
	fullPath := filepath.Join(paths...)
	_, err := os.Stat(fullPath)
	exists := !errors.Is(err, os.ErrNotExist)
	return exists
}
