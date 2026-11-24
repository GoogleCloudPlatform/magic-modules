/*
* Copyright 2024 Google LLC. All Rights Reserved.
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
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRunnerCopy(t *testing.T) {
	baseDir := t.TempDir()
	srcDir := filepath.Join(baseDir, "src")
	dstDir := filepath.Join(baseDir, "dst")

	if err := os.MkdirAll(filepath.Join(srcDir, "sub", "deep"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("updated-a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "sub", "deep", "run.sh"), []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}

	// Pre-populate destination to verify overwriting
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dstDir, "a.txt"), []byte("initial-a"), 0644); err != nil {
		t.Fatal(err)
	}

	rnr, err := NewRunner()
	if err != nil {
		t.Fatal(err)
	}
	if err := rnr.PushDir(baseDir); err != nil {
		t.Fatal(err)
	}

	// Copy using relative paths to test Runner.abs resolution too
	if err := rnr.Copy("src", "dst"); err != nil {
		t.Fatalf("Runner.Copy failed: %v", err)
	}

	gotA, err := os.ReadFile(filepath.Join(dstDir, "a.txt"))
	if err != nil {
		t.Fatalf("reading a.txt: %v", err)
	}
	if string(gotA) != "updated-a" {
		t.Errorf("a.txt = %q, want %q", string(gotA), "updated-a")
	}

	runPath := filepath.Join(dstDir, "sub", "deep", "run.sh")
	gotRun, err := os.ReadFile(runPath)
	if err != nil {
		t.Fatalf("reading run.sh: %v", err)
	}
	if string(gotRun) != "#!/bin/sh\n" {
		t.Errorf("run.sh = %q, want %q", string(gotRun), "#!/bin/sh\n")
	}
	info, err := os.Stat(runPath)
	if err != nil {
		t.Fatalf("stat run.sh: %v", err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("run.sh mode = %o, want 0755", info.Mode().Perm())
	}

	if err := rnr.Copy("non-existent-src", "dst"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist for missing source, got %v", err)
	}
}
