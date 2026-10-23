// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldAddHashicorpCopyRightHeader(t *testing.T) {
	cases := []struct {
		name string
		// dir is the output folder, relative to a temporary root.
		dir string
		// module is written to dir/go.mod; no go.mod is written when empty.
		module string
		target string
		want   bool
	}{
		{
			name:   "ga repo folder",
			dir:    "terraform-provider-google",
			module: "github.com/hashicorp/terraform-provider-google",
			target: "google/services/compute/resource_compute_instance.go",
			want:   true,
		},
		{
			name:   "beta worktree named after a task",
			dir:    "worktrees/terraform-provider-google-beta/my-task",
			module: "github.com/hashicorp/terraform-provider-google-beta",
			target: "google-beta/services/compute/resource_compute_instance.go",
			want:   true,
		},
		{
			name:   "nightly worktree named after a task",
			dir:    "worktrees/terraform-provider-google-nightly/my-task",
			module: "github.com/hashicorp/terraform-provider-google-nightly",
			target: "google-nightly/provider/provider.go",
			want:   true,
		},
		{
			name:   "go.mod for a repo without headers wins over folder name",
			dir:    "terraform-provider-google",
			module: "github.com/hashicorp/terraform-provider-google-private",
			target: "google-private/provider/provider.go",
			want:   false,
		},
		{
			name:   "tgc repo",
			dir:    "terraform-google-conversion",
			module: "github.com/GoogleCloudPlatform/terraform-google-conversion/v7",
			target: "pkg/tfplan2cai/converters/google/resources/services/compute/compute_instance.go",
			want:   false,
		},
		{
			name:   "no go.mod falls back to folder name",
			dir:    "terraform-provider-google-beta",
			target: "google-beta/provider/provider.go",
			want:   true,
		},
		{
			name:   "no go.mod and unknown folder name",
			dir:    "worktrees/terraform-provider-google-beta/my-task",
			target: "google-beta/provider/provider.go",
			want:   false,
		},
		{
			name:   "ignored file in a header repo",
			dir:    "worktrees/terraform-provider-google/my-task",
			module: "github.com/hashicorp/terraform-provider-google",
			target: "scripts/sweeper.go",
			want:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outputFolder := filepath.Join(t.TempDir(), tc.dir)
			if err := os.MkdirAll(outputFolder, 0755); err != nil {
				t.Fatal(err)
			}
			if tc.module != "" {
				goMod := "module " + tc.module + "\n\ngo 1.26.0\n"
				if err := os.WriteFile(filepath.Join(outputFolder, "go.mod"), []byte(goMod), 0644); err != nil {
					t.Fatal(err)
				}
			}

			if got := (Terraform{}).shouldAddHashicorpCopyRightHeader(outputFolder, tc.target); got != tc.want {
				t.Errorf("shouldAddHashicorpCopyRightHeader(%q, %q) = %v, want %v", tc.dir, tc.target, got, tc.want)
			}
		})
	}
}

func TestOutputModulePathIsCached(t *testing.T) {
	outputFolder := t.TempDir()
	goModPath := filepath.Join(outputFolder, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module github.com/hashicorp/terraform-provider-google-beta\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if got, want := outputModulePath(outputFolder), "github.com/hashicorp/terraform-provider-google-beta"; got != want {
		t.Fatalf("outputModulePath() = %q, want %q", got, want)
	}

	// Generation rewrites go.mod mid-run; later lookups must keep the first answer.
	if err := os.WriteFile(goModPath, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if got, want := outputModulePath(outputFolder), "github.com/hashicorp/terraform-provider-google-beta"; got != want {
		t.Errorf("outputModulePath() after rewrite = %q, want %q", got, want)
	}
}
