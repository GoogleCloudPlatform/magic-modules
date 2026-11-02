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
package github

import (
	"testing"

	gh "github.com/google/go-github/v68/github"
)

func TestConvertGHPullRequestHeadSha(t *testing.T) {
	cases := map[string]struct {
		pr   *gh.PullRequest
		want string
	}{
		"head sha": {
			pr:   &gh.PullRequest{Head: &gh.PullRequestBranch{SHA: gh.Ptr("abc123")}},
			want: "abc123",
		},
		"no head": {
			pr:   &gh.PullRequest{},
			want: "",
		},
		"nil pull request": {
			pr:   nil,
			want: "",
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if got := convertGHPullRequest(tc.pr).HeadSha; got != tc.want {
				t.Errorf("convertGHPullRequest().HeadSha = %q, want %q", got, tc.want)
			}
		})
	}
}
