/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
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

	ghi "github.com/google/go-github/v68/github"
)

func TestExecRemoveAwaitingApproval(t *testing.T) {
	mockGH := &mockGithub{
		pullRequest: &ghi.PullRequest{
			User: &ghi.User{
				Login: ghi.String("core_author"),
			},
		},
		userType:      github.CoreContributorUserType,
		calledMethods: make(map[string][][]any),
	}

	execRemoveLabel("pr1", mockGH, "awaiting-approval")

	method := "RemoveLabel"
	expected := [][]any{{"pr1", "awaiting-approval"}}

	if calls, ok := mockGH.calledMethods[method]; !ok {
		t.Fatal("awaiting-approval label not removed for PR ")
	} else if !reflect.DeepEqual(calls, expected) {
		t.Fatalf("Wrong calls for %s, got %v, expected %v", method, calls, expected)
	}
}
