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
	"magician/github"
	"reflect"
	"testing"
)

func TestExecRemoveAwaitingApproval(t *testing.T) {
	gh := &mockGithub{
		pullRequest: github.PullRequest{
			User: github.User{
				Login: "core_author",
			},
		},
		userType:      github.CoreContributorUserType,
		calledMethods: make(map[string][][]any),
	}

	execRemoveLabel("pr1", gh, "awaiting-approval")

	method := "RemoveLabel"
	expected := [][]any{{"pr1", "awaiting-approval"}}
	if calls, ok := gh.calledMethods[method]; !ok {
		t.Fatal("awaiting-approval label not removed for PR ")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}

}
