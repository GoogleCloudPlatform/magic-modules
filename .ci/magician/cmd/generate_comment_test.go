package cmd

import (
	"reflect"
	"testing"
)

func TestExecGenerateComment(t *testing.T) {
	r := &mockRunner{
		calledMethods: make(map[string][][]any),
		cmdResults: map[string]string{
			"/mock/dir/tfc git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/docs-examples /mock/dir/tfoics] []":                "",
			"/mock/dir/tpgb git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion /mock/dir/tfc] []":    "",
			" git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/terraform-provider-google /mock/dir/tpg] []":                    "",
			"/mock/dir/tpg git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta /mock/dir/tpgb] []": "",
			"/mock/dir/magic-modules git [diff HEAD origin/main tools/missing-test-detector] []":                                                                     "",
			"/mock/dir/magic-modules/tools/diff-processor bin/diff-processor [breaking-changes] []":                                                                  "",
			"/mock/dir/magic-modules/tools/diff-processor make [build] [OLD_REF=auto-pr-pr1-old NEW_REF=auto-pr-pr1]":                                                "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [mod edit -replace google/provider/new=/mock/dir/tpgb] []":                                       "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [mod edit -replace google/provider/old=/mock/dir/tpgbold] []":                                    "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [mod tidy] []":                                                                                   "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [run . -services-dir=/mock/dir/tpgb/google-beta/services] []":                                    "## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n",
			"/mock/dir/tfc git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                      " 1 file changed, 10 insertions(+)\n",
			"/mock/dir/tfc git [fetch origin auto-pr-pr1-old] []":                                                                                                    "",
			"/mock/dir/tfoics git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                   "",
			"/mock/dir/tfoics git [fetch origin auto-pr-pr1-old] []":                                                                                                 "",
			"/mock/dir/tpg git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                      " 2 files changed, 40 insertions(+)\n",
			"/mock/dir/tpg git [fetch origin auto-pr-pr1-old] []":                                                                                                    "",
			"/mock/dir/tpgb find [. -type f -name *.go -exec sed -i.bak s~github.com/hashicorp/terraform-provider-google-beta~google/provider/new~g {} +] []":        "",
			"/mock/dir/tpgb git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                     " 2 files changed, 40 insertions(+)\n",
			"/mock/dir/tpgb git [fetch origin auto-pr-pr1-old] []":                                                                                                   "",
			"/mock/dir/tpgb sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g go.mod] []":                                      "",
			"/mock/dir/tpgb sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g go.sum] []":                                      "",
			"/mock/dir/tpgbold find [. -type f -name *.go -exec sed -i.bak s~github.com/hashicorp/terraform-provider-google-beta~google/provider/old~g {} +] []":     "",
			"/mock/dir/tpgbold git [checkout origin/auto-pr-pr1-old] []":                                                                                             "",
			"/mock/dir/tpgbold sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g go.mod] []":                                   "",
			"/mock/dir/tpgbold sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g go.sum] []":                                   "",
		},
	}
	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	execGenerateComment("build1", "project1", "17", "sha1", "pr1", "*******", gh, r)

	for method, expectedCalls := range map[string][][]any{
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
		},
		"Run": {
			{"", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google", "/mock/dir/tpg"}, []string(nil)},
			{"/mock/dir/tpg", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, []string(nil)},
			{"/mock/dir/tpg", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, []string(nil)},
			{"/mock/dir/tpg", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta", "/mock/dir/tpgb"}, []string(nil)},
			{"/mock/dir/tpgb", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, []string(nil)},
			{"/mock/dir/tpgb", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, []string(nil)},
			{"/mock/dir/tpgb", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion", "/mock/dir/tfc"}, []string(nil)},
			{"/mock/dir/tfc", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, []string(nil)},
			{"/mock/dir/tfc", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, []string(nil)},
			{"/mock/dir/tfc", "git", []string{"clone", "-b", "auto-pr-pr1", "https://modular-magician:*******@github.com/modular-magician/docs-examples", "/mock/dir/tfoics"}, []string(nil)},
			{"/mock/dir/tfoics", "git", []string{"fetch", "origin", "auto-pr-pr1-old"}, []string(nil)},
			{"/mock/dir/tfoics", "git", []string{"diff", "origin/auto-pr-pr1-old", "origin/auto-pr-pr1", "--shortstat"}, []string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, []string{"OLD_REF=auto-pr-pr1-old", "NEW_REF=auto-pr-pr1"}},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, []string(nil)},
			{"/mock/dir/tpgbold", "git", []string{"checkout", "origin/auto-pr-pr1-old"}, []string(nil)},
			{"/mock/dir/tpgbold", "find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", "s~github.com/hashicorp/terraform-provider-google-beta~google/provider/old~g", "{}", "+"}, []string(nil)},
			{"/mock/dir/tpgbold", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g", "go.mod"}, []string(nil)},
			{"/mock/dir/tpgbold", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g", "go.sum"}, []string(nil)},
			{"/mock/dir/tpgb", "find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", "s~github.com/hashicorp/terraform-provider-google-beta~google/provider/new~g", "{}", "+"}, []string(nil)},
			{"/mock/dir/tpgb", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g", "go.mod"}, []string(nil)},
			{"/mock/dir/tpgb", "sed", []string{"-i.bak", "s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g", "go.sum"}, []string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"mod", "edit", "-replace", "google/provider/new=/mock/dir/tpgb"}, []string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"mod", "edit", "-replace", "google/provider/old=/mock/dir/tpgbold"}, []string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"mod", "tidy"}, []string(nil)},
			{"/mock/dir/magic-modules/tools/missing-test-detector", "go", []string{"run", ".", "-services-dir=/mock/dir/tpgb/google-beta/services"}, []string(nil)},
			{"/mock/dir/magic-modules", "git", []string{"diff", "HEAD", "origin/main", "tools/missing-test-detector"}, []string(nil)}},
	} {
		if actualCalls, ok := r.calledMethods[method]; !ok {
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
		"PostComment":     {{"pr1", "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n## Diff report\nYour PR generated some diffs in downstreams - here they are.\n\nTerraform GA: [Diff](https://github.com/modular-magician/terraform-provider-google/compare/auto-pr-pr1-old..auto-pr-pr1) ( 2 files changed, 40 insertions(+))\nTerraform Beta: [Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/auto-pr-pr1-old..auto-pr-pr1) ( 2 files changed, 40 insertions(+))\nTF Conversion: [Diff](https://github.com/modular-magician/terraform-google-conversion/compare/auto-pr-pr1-old..auto-pr-pr1) ( 1 file changed, 10 insertions(+))\n## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n\n"}},
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
