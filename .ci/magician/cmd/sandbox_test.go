/*
* Copyright 2026 Google LLC. All Rights Reserved.
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
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"magician/exec"
)

type sandbox struct {
	Dir    string
	Runner ExecRunner
}

// Routes some commands to fake sandbox scripts, bypassing $PATH.
type interceptingRunner struct {
	*exec.Runner
	sbDir              string
	allowedPassthrough map[string]bool
	strictMode         bool
}

func (s *interceptingRunner) Run(name string, args []string, env map[string]string) (string, error) {
	if !s.strictMode {
		if !strings.Contains(name, string(filepath.Separator)) {
			sandboxBin := filepath.Join(s.sbDir, name)
			if _, err := os.Stat(sandboxBin); err == nil {
				name = sandboxBin
			}
		}
		return s.Runner.Run(name, args, env)
	}

	baseName := filepath.Base(name)

	sandboxBin := filepath.Join(s.sbDir, baseName)
	if _, err := os.Stat(sandboxBin); err == nil {
		return s.Runner.Run(sandboxBin, args, env)
	}

	if s.allowedPassthrough[baseName] {
		return s.Runner.Run(name, args, env)
	}

	return "", fmt.Errorf("command %q is not in the passthrough list and no sandbox mock was found. Please mock it or add it to the passthrough list", baseName)
}

func newSandbox(t *testing.T) *sandbox {
	dir := t.TempDir()

	realRunner, err := exec.NewRunner()
	if err != nil {
		t.Fatalf("Failed to create runner: %v", err)
	}

	if err := realRunner.PushDir(dir); err != nil {
		t.Fatalf("Failed to push dir: %v", err)
	}

	realRunner.MustRun("git", []string{"init", "-b", "main"}, nil)

	realRunner.MustRun("touch", []string{"README.md"}, nil)
	realRunner.MustRun("git", []string{"add", "."}, nil)
	realRunner.MustRun("git", []string{"config", "user.email", "test@example.com"}, nil)
	realRunner.MustRun("git", []string{"config", "user.name", "Test Sandbox"}, nil)
	realRunner.MustRun("git", []string{"config", "commit.gpgsign", "false"}, nil)
	realRunner.MustRun("git", []string{"commit", "-m", "Initial sandbox commit"}, nil)

	runner := &interceptingRunner{
		Runner: realRunner,
		sbDir:  dir,
		allowedPassthrough: map[string]bool{
			"chmod": true,
			"rm":    true,
			"mkdir": true,
			"touch": true,
			"cp":    true,
			"mv":    true,
			"ls":    true,
			"echo":  true,
			"cat":   true,
			"grep":  true,
		},
	}

	return &sandbox{
		Dir:    dir,
		Runner: runner,
	}
}

func (s *sandbox) AllowPassthrough(commands ...string) {
	if runner, ok := s.Runner.(*interceptingRunner); ok {
		for _, cmd := range commands {
			runner.allowedPassthrough[cmd] = true
		}
	}
}

func (s *sandbox) RequireAllowlist() {
	if runner, ok := s.Runner.(*interceptingRunner); ok {
		runner.strictMode = true
	}
}
