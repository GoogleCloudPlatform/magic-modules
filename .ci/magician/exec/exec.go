package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	cp "github.com/otiai10/copy"
)

type actualRunner struct{}

func (ar *actualRunner) Getwd() (string, error) {
	return os.Getwd()
}

func (ar *actualRunner) Copy(src, dest string) error {
	return cp.Copy(src, dest)
}

func (ar *actualRunner) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ar actualRunner) Chdir(path string) {
	os.Chdir(path)
}

func (ar actualRunner) Run(name string, args, env []string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.Output()
	if err != nil {
		exitErr := err.(*exec.ExitError)
		return string(out), fmt.Errorf("error running %s: %v\nstdout:\n%sstderr:\n%s", name, err, out, exitErr.Stderr)
	}
	return string(out), nil
}

func (ar actualRunner) MustRun(name string, args, env []string) string {
	out, err := ar.Run(name, args, env)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

type Runner interface {
	Getwd() (string, error)
	Copy(src, dest string) error
	RemoveAll(path string) error
	Chdir(path string)
	Run(name string, args, env []string) (string, error)
	MustRun(name string, args, env []string) string
}

func NewRunner() Runner {
	return &actualRunner{}
}
