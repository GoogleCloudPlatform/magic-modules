package exec

import (
	"container/list"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

type actualRunner struct {
	cwd      string
	dirStack *list.List
}

type Runner interface {
	GetCWD() string
	Copy(src, dest string) error
	Mkdir(path string) error
	Walk(root string, fn WalkFunc) error
	RemoveAll(path string) error
	PushDir(path string) error
	PopDir() error
	WriteFile(name, data string) error
	ReadFile(name string) (string, error)
	Run(name string, args, env []string) (string, error)
	MustRun(name string, args, env []string) string
}

func NewRunner() (Runner, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &actualRunner{
		cwd:      wd,
		dirStack: list.New(),
	}, nil
}

func (ar *actualRunner) GetCWD() string {
	return ar.cwd
}

func (ar *actualRunner) Copy(src, dest string) error {
	return cp.Copy(ar.abs(src), ar.abs(dest))
}

func (ar *actualRunner) Mkdir(path string) error {
	return os.MkdirAll(ar.abs(path), 0777)
}

func (ar *actualRunner) Walk(root string, fn WalkFunc) error {
	return filepath.Walk(root, fn)
}

func (ar *actualRunner) RemoveAll(path string) error {
	return os.RemoveAll(ar.abs(path))
}

// PushDir changes the directory for the runner to the desired path and saves the previous directory in the stack.
func (ar *actualRunner) PushDir(path string) error {
	if ar.dirStack == nil {
		return errors.New("attempted to push dir, but stack was nil")
	}
	ar.dirStack.PushFront(ar.cwd)
	ar.cwd = ar.abs(path)
	return nil
}

// PopDir removes the most recently added directory from the stack and changes front to it.
func (ar *actualRunner) PopDir() error {
	if ar.dirStack == nil || ar.dirStack.Len() == 0 {
		return errors.New("attempted to pop dir, but stack was nil or empty")
	}
	frontVal := ar.dirStack.Remove(ar.dirStack.Front())
	dir, ok := frontVal.(string)
	if !ok {
		return fmt.Errorf("last element in dir stack was a %T, expected string", frontVal)
	}
	ar.cwd = dir
	return nil
}

func (ar *actualRunner) WriteFile(name, data string) error {
	return os.WriteFile(ar.abs(name), []byte(data), 0644)
}

func (ar *actualRunner) ReadFile(name string) (string, error) {
	data, err := os.ReadFile(ar.abs(name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Run the given command with the given args and env, return output and error if any
func (ar *actualRunner) Run(name string, args, env []string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = ar.cwd
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.Output()
	switch typedErr := err.(type) {
	case *exec.ExitError:
		return string(out), fmt.Errorf("error running %s: %v\nstdout:\n%sstderr:\n%s", name, err, out, typedErr.Stderr)
	case *fs.PathError:
		return "", fmt.Errorf("path error running %s: %v", name, typedErr)

	}
	return string(out), nil
}

// Run the command and exit if there's an error.
func (ar *actualRunner) MustRun(name string, args, env []string) string {
	out, err := ar.Run(name, args, env)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func (ar *actualRunner) abs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(ar.cwd, path)
	}
	return path
}
