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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSandboxInitialization(t *testing.T) {
	sb := newSandbox(t)

	if _, err := os.Stat(sb.Dir); os.IsNotExist(err) {
		t.Fatalf("Sandbox directory %s does not exist", sb.Dir)
	}

	sb.AllowPassthrough("git")
	output, err := sb.Runner.Run("git", []string{"status"}, nil)
	if err != nil {
		t.Fatalf("Sandbox is not a valid git repository: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "On branch main") {
		t.Errorf("Expected 'On branch main' in git status output, got: %s", output)
	}
}

func TestSandboxInterceptor_AllowedCommand(t *testing.T) {
	sb := newSandbox(t)
	sb.RequireAllowlist()
	_, err := sb.Runner.Run("touch", []string{"test_file.txt"}, nil)
	if err != nil {
		t.Fatalf("Expected allowlisted command to succeed, got: %v", err)
	}
}

func TestSandboxInterceptor_MockOverridesAllowed(t *testing.T) {
	sb := newSandbox(t)
	sb.RequireAllowlist()
	sb.Runner.WriteFile("echo", "#!/bin/bash\necho 'intercepted echo'")
	os.Chmod(filepath.Join(sb.Dir, "echo"), 0755)

	output, err := sb.Runner.Run("echo", []string{"hello"}, nil)
	if err != nil {
		t.Fatalf("Expected mocked command to succeed, got: %v", err)
	}
	if !strings.Contains(output, "intercepted echo") {
		t.Errorf("Expected sandbox mock to override allowlist, got: %s", output)
	}
}

func TestSandboxInterceptor_UnhandledCommandPanics(t *testing.T) {
	sb := newSandbox(t)
	sb.RequireAllowlist()
	_, err := sb.Runner.Run("curl", []string{"http://example.com"}, nil)
	if err == nil {
		t.Fatalf("Expected unhandled command to return error, but it succeeded")
	}
	if !strings.Contains(err.Error(), "not in the passthrough list") {
		t.Errorf("Expected strict interceptor error, got: %v", err)
	}
}
