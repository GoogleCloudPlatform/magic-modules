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
	"testing"

	"magician/exec"
)

type sandbox struct {
	Dir    string
	Runner ExecRunner
}

func newSandbox(t *testing.T) *sandbox {
	dir := t.TempDir()

	runner, err := exec.NewRunner()
	if err != nil {
		t.Fatalf("Failed to create runner: %v", err)
	}

	if err := runner.PushDir(dir); err != nil {
		t.Fatalf("Failed to push dir: %v", err)
	}

	runner.MustRun("git", []string{"init", "-b", "main"}, nil)

	runner.MustRun("touch", []string{"README.md"}, nil)
	runner.MustRun("git", []string{"add", "."}, nil)
	runner.MustRun("git", []string{"config", "user.email", "test@example.com"}, nil)
	runner.MustRun("git", []string{"config", "user.name", "Test Sandbox"}, nil)
	runner.MustRun("git", []string{"config", "commit.gpgsign", "false"}, nil)
	runner.MustRun("git", []string{"commit", "-m", "Initial sandbox commit"}, nil)

	return &sandbox{
		Dir:    dir,
		Runner: runner,
	}
}
