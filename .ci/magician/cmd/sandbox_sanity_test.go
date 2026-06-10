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
	"strings"
	"testing"
)

func TestSandboxInitialization(t *testing.T) {
	sb := newSandbox(t)

	if _, err := os.Stat(sb.Dir); os.IsNotExist(err) {
		t.Fatalf("Sandbox directory %s does not exist", sb.Dir)
	}

	output, err := sb.Runner.Run("git", []string{"status"}, nil)
	if err != nil {
		t.Fatalf("Sandbox is not a valid git repository: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "On branch main") {
		t.Errorf("Expected 'On branch main' in git status output, got: %s", output)
	}
}
