package verify

import (
	"errors"
	"os"
	"path/filepath"
)

// If the argument is a path, pathOrContents loads it and returns the contents,
// otherwise the argument is assumed to be the desired contents and is simply
// returned.
//
// The boolean second return value can be called `wasPath` - it indicates if a
// path was detected and a file loaded.
func PathOrContents(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := poc
	if path[0] == '~' {
		if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
			return path, true, errors.New("cannot expand user-specific home dir")
		}
		dir, err := os.UserHomeDir()
		if err != nil {
			return path, true, err
		}
		path = filepath.Join(dir, path[1:])
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}
