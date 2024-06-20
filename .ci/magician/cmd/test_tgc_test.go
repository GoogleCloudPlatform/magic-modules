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
package cmd

import (
	"reflect"
	"testing"
)

func TestExecTestTGC(t *testing.T) {
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}

	execTestTGC("sha1", "pr1", gh)

	method := "CreateWorkflowDispatchEvent"
	expected := [][]any{{"test-tgc.yml", map[string]any{"branch": "auto-pr-pr1", "owner": "modular-magician", "repo": "terraform-google-conversion", "sha": "sha1"}}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("Workflow dispatch event not created")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}
