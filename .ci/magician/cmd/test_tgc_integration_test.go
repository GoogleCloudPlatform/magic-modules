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
	"magician/source"
	"reflect"
	"testing"
)

func TestShouldRunTests(t *testing.T) {
	cases := []struct {
		name         string
		changedFiles []string
		expected     bool
	}{
		{
			name:         "relevant go file",
			changedFiles: []string{"test/services/alloydb/alloydb_cluster_generated_test.go"},
			expected:     true,
		},
		{
			name:         "non-go file",
			changedFiles: []string{"docs/supported_resources.md"},
			expected:     false,
		},
		{
			name:         "ignored directory cai2hcl",
			changedFiles: []string{"cai2hcl/services/certificatemanager/certificate.go"},
			expected:     false,
		},
		{
			name:         "ignored directory tfplan2cai",
			changedFiles: []string{"tfplan2cai/converters/google/resources/services/biglakeiceberg"},
			expected:     false,
		},
		{
			name:         "pkg/services file (ignored by default)",
			changedFiles: []string{"pkg/services/compute/compute_disk.go"},
			expected:     false,
		},
		{
			name:         "pkg/provider/provider_register_services.go file (ignored by default)",
			changedFiles: []string{"pkg/provider/provider_register_services.go"},
			expected:     false,
		},
		{
			name:         "pkg/services cai2hcl file (no longer exception)",
			changedFiles: []string{"pkg/services/compute/compute_disk_cai2hcl.go"},
			expected:     false,
		},
		{
			name:         "pkg/services tfplan2cai file (no longer exception)",
			changedFiles: []string{"pkg/services/compute/compute_disk_tfplan2cai.go"},
			expected:     false,
		},
		{
			name:         "pkg/cai2hcl file",
			changedFiles: []string{"pkg/cai2hcl/converters/convert_resource.go"},
			expected:     true,
		},
		{
			name:         "pkg/tfplan2cai file",
			changedFiles: []string{"pkg/tfplan2cai/converters/cai/convert.go"},
			expected:     true,
		},
		{
			name:         "pkg/caiasset file",
			changedFiles: []string{"pkg/caiasset/asset.go"},
			expected:     true,
		},
		{
			name:         "multiple files, one relevant",
			changedFiles: []string{"README.md", "caiasset/asset.go", "pkg/cai2hcl/converters/convert_resource.go"},
			expected:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := shouldRunTests(tc.changedFiles)
			if actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestExecTestTGCIntegration_DiffsAgainstOldBranch(t *testing.T) {
	mr := NewMockRunner().(*mockRunner)
	repoPath := "/go/src/github.com/modular-magician/terraform-google-conversion"
	mr.cmdResults[repoPath+" git [diff origin/auto-pr-12345-old origin/auto-pr-12345 --name-only] map[]"] = "docs/supported_resources.md\n"
	ctlr := source.NewController("/go", "modular-magician", "secret-token", mr)
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}

	err := execTestTGCIntegration("12345", "sha123", "build1", "proj1", "19", "terraform-google-conversion", "modular-magician", mr, ctlr, gh)
	if err != nil {
		t.Fatalf("execTestTGCIntegration() unexpected error: %v", err)
	}

	calls, ok := mr.Calls("Run")
	if !ok {
		t.Fatal("expected Run calls on mockRunner")
	}
	wantCalls := []ParameterList{
		{
			"/mock/dir/magic-modules/.ci/magician",
			"git",
			[]string{"clone", "-b", "auto-pr-12345", "https://modular-magician:secret-token@github.com/modular-magician/terraform-google-conversion", repoPath},
			map[string]string(nil),
		},
		{
			repoPath,
			"git",
			[]string{"fetch", "origin", "auto-pr-12345-old"},
			map[string]string(nil),
		},
		{
			repoPath,
			"git",
			[]string{"diff", "origin/auto-pr-12345-old", "origin/auto-pr-12345", "--name-only"},
			map[string]string(nil),
		},
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Errorf("Run calls = %v, want %v", calls, wantCalls)
	}
}
