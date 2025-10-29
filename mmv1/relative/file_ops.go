package relative

import (
	"fmt"
	"os"
	"path/filepath"
)

// For file IO relative to some base directory
// Templates and other files are accessed
// relative to os CWD, this library allows accessing a
// location relative to an initialized global base

var baseDirectory string

func ReadFile(path string) ([]byte, error) {
	fmt.Printf("basedir is : `%s`", baseDirectory)
	absolutePath := filepath.Join(baseDirectory, path)
	return os.ReadFile(absolutePath)
}

func SetBaseDir(dir string) {
	baseDirectory = dir
	fmt.Printf("set base dir to : `%s`", dir)
}
