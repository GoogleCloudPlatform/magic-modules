package google

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/golang/glog"
)

var ModuleRoot = ""

func GetModuleRoot() string {
	if ModuleRoot != "" {
		return ModuleRoot
	}

	// First get the current file's path
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)

	// Walk up the directory tree until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			ModuleRoot = dir
			return ModuleRoot
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root without finding go.mod
			glog.Exitf("Could not find module root (no go.mod file found)\nStack trace:\n%s", debug.Stack())
		}
		dir = parent
	}
}
