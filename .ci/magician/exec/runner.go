/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
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

type Runner struct {
	cwd      string
	dirStack *list.List
}

func NewRunner() (*Runner, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Runner{
		cwd:      wd,
		dirStack: list.New(),
	}, nil
}

func (ar *Runner) GetCWD() string {
	return ar.cwd
}

func (ar *Runner) Copy(src, dest string) error {
	return cp.Copy(ar.abs(src), ar.abs(dest))
}

func (ar *Runner) Mkdir(path string) error {
	return os.MkdirAll(ar.abs(path), 0777)
}

func (ar *Runner) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

func (ar *Runner) RemoveAll(path string) error {
	return os.RemoveAll(ar.abs(path))
}

// PushDir changes the directory for the runner to the desired path and saves the previous directory in the stack.
func (ar *Runner) PushDir(path string) error {
	if ar.dirStack == nil {
		return errors.New("attempted to push dir, but stack was nil")
	}
	ar.dirStack.PushFront(ar.cwd)
	ar.cwd = ar.abs(path)
	return nil
}

// PopDir removes the most recently added directory from the stack and changes front to it.
func (ar *Runner) PopDir() error {
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

func (ar *Runner) WriteFile(name, data string) error {
	return os.WriteFile(ar.abs(name), []byte(data), 0644)
}

func (ar *Runner) ReadFile(name string) (string, error) {
	data, err := os.ReadFile(ar.abs(name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Run the given command with the given args and env, return output and error if any
func (ar *Runner) Run(name string, args []string, env map[string]string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = ar.cwd
	for ev, val := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", ev, val))
	}
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
func (ar *Runner) MustRun(name string, args []string, env map[string]string) string {
	out, err := ar.Run(name, args, env)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func (ar *Runner) abs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(ar.cwd, path)
	}
	return path
}
