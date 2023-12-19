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
	"magician/source"
	"reflect"
	"testing"
)

func TestExecGenerateComment(t *testing.T) {
	mr := NewMockRunner()
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	ctlr := source.NewController("/mock/dir/go", "modular-magician", "*******", mr)
	env := map[string]string{
		"BUILD_ID":     "build1",
		"BUILD_STEP":   "17",
		"COMMIT_SHA":   "sha1",
		"GITHUB_TOKEN": "*******",
		"PR_NUMBER":    "pr1",
		"PROJECT_ID":   "project1",
	}
	diffProcessorEnv := map[string]string{
		"BUILD_ID":     "build1",
		"BUILD_STEP":   "17",
		"COMMIT_SHA":   "sha1",
		"GITHUB_TOKEN": "*******",
		"NEW_REF":      "auto-pr-pr1",
		"OLD_REF":      "auto-pr-pr1-old",
		"PR_NUMBER":    "pr1",
		"PROJECT_ID":   "project1",
	}
	execGenerateComment(env, gh, mr, ctlr)

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
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google", "/mock/dir/tpg"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta", "/mock/dir/tpgb"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion", "/mock/dir/tfc"}, map[string]string(nil)},
			{"/mock/dir/tfc", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, map[string]string(nil)},
			{"/mock/dir/tfc", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/docs-examples", "/mock/dir/tfoics"}, map[string]string(nil)},
			{"/mock/dir/tfoics", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, map[string]string(nil)},
			{"/mock/dir/tfoics", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"add-labels", "pr1"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"add-labels", "pr1"}, diffProcessorEnv},
			{"/mock/dir/tpgbold", "git", []string{"checkout", "origin/auto-pr-pr1-old"}, map[string]string(nil)},
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
		"PostBuildStatus": {{"pr1", "terraform-provider-breaking-change-test", "success", "https://console.cloud.google.com/cloud-build/builds;region=global/build1;step=17?project=project1", "sha1"}},
		"PostComment":     {{"pr1", "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n## Diff report\nYour PR generated some diffs in downstreams - here they are.\n\nTerraform GA: [Diff](https://github.com/modular-magician/terraform-provider-google/compare/auto-pr-pr1-old..auto-pr-pr1) ( 2 files changed, 40 insertions(+))\nTerraform Beta: [Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/auto-pr-pr1-old..auto-pr-pr1) ( 2 files changed, 40 insertions(+))\nTF Conversion: [Diff](https://github.com/modular-magician/terraform-google-conversion/compare/auto-pr-pr1-old..auto-pr-pr1) ( 1 file changed, 10 insertions(+))\n\n## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n\n"}},
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
