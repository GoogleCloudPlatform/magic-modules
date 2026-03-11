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
)

func TestShouldRunTests(t *testing.T) {
	cases := []struct {
		name         string
		changedFiles []string
		expected     bool
	}{
		{
			name:         "relevant go file",
			changedFiles: []string{"mmv1/provider/terraform_tgc.go"},
			expected:     true,
		},
		{
			name:         "non-go file",
			changedFiles: []string{"mmv1/provider/terraform_tgc.yaml"},
			expected:     false,
		},
		{
			name:         "ignored directory cai2hcl",
			changedFiles: []string{"cai2hcl/main.go"},
			expected:     false,
		},
		{
			name:         "ignored directory tfplan2cai",
			changedFiles: []string{"tfplan2cai/main.go"},
			expected:     false,
		},
		{
			name:         "pkg/services file (ignored by default)",
			changedFiles: []string{"pkg/services/compute/compute.go"},
			expected:     false,
		},
		{
			name:         "pkg/services cai2hcl file (exception)",
			changedFiles: []string{"pkg/services/compute/compute_cai2hcl.go"},
			expected:     true,
		},
		{
			name:         "pkg/services tfplan2cai file (exception)",
			changedFiles: []string{"pkg/services/compute/tfplan2cai.go"},
			expected:     true,
		},
		{
			name:         "multiple files, one relevant",
			changedFiles: []string{"README.md", "cai2hcl/main.go", "mmv1/provider/terraform_tgc.go"},
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
