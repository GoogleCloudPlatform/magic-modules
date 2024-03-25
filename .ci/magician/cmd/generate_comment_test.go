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
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"magician/source"
)

func TestExecGenerateComment(t *testing.T) {
	mr := NewMockRunner()
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	ctlr := source.NewController("/mock/dir/go", "modular-magician", "*******", mr)
	diffProcessorEnv := map[string]string{
		"NEW_REF": "auto-pr-123456",
		"OLD_REF": "auto-pr-123456-old",
		"PATH":    os.Getenv("PATH"),
		"GOPATH":  os.Getenv("GOPATH"),
		"HOME":    os.Getenv("HOME"),
	}
	execGenerateComment(
		123456,
		"*******",
		"build1",
		"17",
		"project1",
		"sha1",
		gh,
		mr,
		ctlr,
	)

	for method, expectedCalls := range map[string][]ParameterList{
		"Copy": {
			{"/mock/dir/tpg", "/mock/dir/magic-modules/tools/diff-processor/old"},
			{"/mock/dir/tpg", "/mock/dir/magic-modules/tools/diff-processor/new"},
			{"/mock/dir/tpgb", "/mock/dir/magic-modules/tools/diff-processor/old"},
			{"/mock/dir/tpgb", "/mock/dir/magic-modules/tools/diff-processor/new"},
			{"/mock/dir/tpgb", "/mock/dir/tpgbold"},
		},
		"RemoveAll": {
			{"/mock/dir/magic-modules/tools/diff-processor/old"},
			{"/mock/dir/magic-modules/tools/diff-processor/new"},
			{"/mock/dir/magic-modules/tools/diff-processor/bin"},
			{"/mock/dir/magic-modules/tools/diff-processor/old"},
			{"/mock/dir/magic-modules/tools/diff-processor/new"},
			{"/mock/dir/magic-modules/tools/diff-processor/bin"},
		},
		"Run": {
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google", "/mock/dir/tpg"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta", "/mock/dir/tpgb"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion", "/mock/dir/tgc"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/docs-examples", "/mock/dir/tfoics"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/tgc", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tgc", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/tfoics", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tfoics", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"changed-schema-labels"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"changed-schema-labels"}, map[string]string(nil)},
			{"/mock/dir/tpgbold", "git", []string{"checkout", "origin/auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpgbold", "find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", "s~github.com/hashicorp/terraform-provider-google-beta~google/provider/old~g", "{}", "+"}, map[string]string(nil)},
			{"/mock/dir/tpgbold", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g", "go.mod"}, map[string]string(nil)},
			{"/mock/dir/tpgbold", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g", "go.sum"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", "s~github.com/hashicorp/terraform-provider-google-beta~google/provider/new~g", "{}", "+"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g", "go.mod"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g", "go.sum"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"mod", "edit", "-replace", "google/provider/new=/mock/dir/tpgb"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"mod", "edit", "-replace", "google/provider/old=/mock/dir/tpgbold"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"mod", "tidy"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"run", ".", "-services-dir=/mock/dir/tpgb/google-beta/services"}, map[string]string(nil)},
			{"/mock/dir/magic-modules", "git", []string{"diff", "HEAD", "origin/main", "tools/missing-test-detector"}, map[string]string(nil)},
		},
	} {
		if actualCalls, ok := mr.Calls(method); !ok {
			t.Fatalf("Found no calls for %s", method)
		} else if len(actualCalls) != len(expectedCalls) {
			t.Fatalf("Unexpected number of calls for %s, got %d, expected %d", method, len(actualCalls), len(expectedCalls))
		} else {
			for i, actualParams := range actualCalls {
				if expectedParams := expectedCalls[i]; !reflect.DeepEqual(actualParams, expectedParams) {
					t.Fatalf("Wrong params for call %d to %s, got %v, expected %v", i, method, actualParams, expectedParams)
				}
			}
		}
	}

	for method, expectedCalls := range map[string][][]any{
		"PostBuildStatus": {{"123456", "terraform-provider-breaking-change-test", "success", "https://console.cloud.google.com/cloud-build/builds;region=global/build1;step=17?project=project1", "sha1"}},
		"PostComment":     {{"123456", "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n## Diff report\n\nYour PR generated some diffs in downstreams - here they are.\n\n`google` provider: [Diff](https://github.com/modular-magician/terraform-provider-google/compare/auto-pr-123456-old..auto-pr-123456) ( 2 files changed, 40 insertions(+))\n`google-beta` provider: [Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/auto-pr-123456-old..auto-pr-123456) ( 2 files changed, 40 insertions(+))\n`terraform-google-conversion`: [Diff](https://github.com/modular-magician/terraform-google-conversion/compare/auto-pr-123456-old..auto-pr-123456) ( 1 file changed, 10 insertions(+))\n\n## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n"}},
		"AddLabels":       {{"123456", []string{"service/google-x"}}},
	} {
		if actualCalls, ok := gh.calledMethods[method]; !ok {
			t.Fatalf("Found no calls for %s", method)
		} else if len(actualCalls) != len(expectedCalls) {
			t.Fatalf("Unexpected number of calls for %s, got %d, expected %d", method, len(actualCalls), len(expectedCalls))
		} else {
			for i, actualParams := range actualCalls {
				if expectedParams := expectedCalls[i]; !reflect.DeepEqual(actualParams, expectedParams) {
					t.Fatalf("Wrong params for call %d to %s, got %v, expected %v", i, method, actualParams, expectedParams)
				}
			}
		}
	}
}

func TestFormatDiffComment(t *testing.T) {
	cases := map[string]struct {
		data               diffCommentData
		expectedStrings    []string
		notExpectedStrings []string
	}{
		"basic message": {
			data:            diffCommentData{},
			expectedStrings: []string{"## Diff report", "hasn't generated any diffs"},
			notExpectedStrings: []string{
				"generated some diffs",
				"## Breaking Change(s) Detected",
				"## Errors",
				"## Missing test report",
			},
		},
		"errors are displayed": {
			data: diffCommentData{
				Errors: []Errors{
					{
						Title:  "`google` provider",
						Errors: []string{"Provider 1"},
					},
					{
						Title:  "Other",
						Errors: []string{"Error 1", "Error 2"},
					},
				},
			},
			expectedStrings: []string{"## Diff report", "## Errors", "`google` provider:\n- Provider 1\n\nOther:\n- Error 1\n- Error 2\n"},
			notExpectedStrings: []string{
				"generated some diffs",
				"## Breaking Change(s) Detected",
				"## Missing test report",
			},
		},
		"diffs are displayed": {
			data: diffCommentData{
				PrNumber: 1234567890,
				Diffs: []Diff{
					{
						Title:     "Repo 1",
						Repo:      "repo-1",
						DiffStats: "+1 added, -1 removed",
					},
					{
						Title:     "Repo 2",
						Repo:      "repo-2",
						DiffStats: "+2 added, -2 removed",
					},
				},
			},
			expectedStrings: []string{
				"## Diff report",
				"generated some diffs",
				"Repo 1: [Diff](https://github.com/modular-magician/repo-1/compare/auto-pr-1234567890-old..auto-pr-1234567890) (+1 added, -1 removed)\nRepo 2: [Diff](https://github.com/modular-magician/repo-2/compare/auto-pr-1234567890-old..auto-pr-1234567890) (+2 added, -2 removed)",
			},
			notExpectedStrings: []string{
				"hasn't generated any diffs",
				"## Breaking Change(s) Detected",
				"## Errors",
				"## Missing test report",
			},
		},
		"breaking changes are displayed": {
			data: diffCommentData{
				BreakingChanges: []string{
					"Breaking change 1",
					"Breaking change 2",
				},
			},
			expectedStrings: []string{
				"## Diff report",
				"## Breaking Change(s) Detected",
				"major release",
				"`override-breaking-change`",
				"- Breaking change 1\n- Breaking change 2\n",
			},
			notExpectedStrings: []string{
				"generated some diffs",
				"## Errors",
				"## Missing test report",
			},
		},
		"missing tests are displayed": {
			data: diffCommentData{
				MissingTests: "## Missing test report",
			},
			expectedStrings: []string{
				"## Diff report",
				"## Missing test report",
			},
			notExpectedStrings: []string{
				"generated some diffs",
				"## Breaking Change(s) Detected",
				"## Errors",
			},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			comment, err := formatDiffComment(tc.data)
			assert.Nil(t, err)

			for _, s := range tc.expectedStrings {
				assert.Contains(t, comment, s)
			}

			for _, s := range tc.notExpectedStrings {
				assert.NotContains(t, comment, s)
			}
		})
	}
}
