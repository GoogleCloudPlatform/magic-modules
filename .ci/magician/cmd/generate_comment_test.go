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

	"magician/source"

	"github.com/stretchr/testify/assert"
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
			{"/mock/dir/tpg", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"checkout", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpg", "make", []string{"build"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"checkout", "auto-pr-123456"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta", "/mock/dir/tpgb"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"checkout", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "make", []string{"build"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"checkout", "auto-pr-123456"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion", "/mock/dir/tgc"}, map[string]string(nil)},
			{"/mock/dir/tgc", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/.ci/magician", "git", []string{"clone", "-b", "auto-pr-123456", "https://modular-magician:*******@github.com/modular-magician/docs-examples", "/mock/dir/tfoics"}, map[string]string(nil)},
			{"/mock/dir/tfoics", "git", []string{"fetch", "origin", "auto-pr-123456-old"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/tpg", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--name-only"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/tpgb", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--name-only"}, map[string]string(nil)},
			{"/mock/dir/tgc", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/tgc", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--name-only"}, map[string]string(nil)},
			{"/mock/dir/tfoics", "git", []string{"diff", "origin/auto-pr-123456-old", "origin/auto-pr-123456", "--shortstat"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"changed-schema-resources"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "make", []string{"build"}, diffProcessorEnv},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"breaking-changes"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"detect-missing-tests", "/mock/dir/tpgb/google-beta/services"}, map[string]string(nil)},
			{"/mock/dir/magic-modules/tools/diff-processor", "bin/diff-processor", []string{"changed-schema-resources"}, map[string]string(nil)},
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
		"PostComment":     {{"123456", "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n## Diff report\n\nYour PR generated some diffs in downstreams - here they are.\n\n`google` provider: [Diff](https://github.com/modular-magician/terraform-provider-google/compare/auto-pr-123456-old..auto-pr-123456) ( 2 files changed, 40 insertions(+))\n`google-beta` provider: [Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/auto-pr-123456-old..auto-pr-123456) ( 2 files changed, 40 insertions(+))\n`terraform-google-conversion`: [Diff](https://github.com/modular-magician/terraform-google-conversion/compare/auto-pr-123456-old..auto-pr-123456) ( 1 file changed, 10 insertions(+))\n\n\n\n## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n"}},
		"AddLabels":       {{"123456", []string{"service/alloydb"}}},
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
						ShortStat: "+1 added, -1 removed",
					},
					{
						Title:     "Repo 2",
						Repo:      "repo-2",
						ShortStat: "+2 added, -2 removed",
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
				BreakingChanges: []BreakingChange{
					{
						Message:                "Breaking change 1",
						DocumentationReference: "doc1",
					},
					{
						Message:                "Breaking change 2",
						DocumentationReference: "doc2",
					},
				},
			},
			expectedStrings: []string{
				"## Diff report",
				"## Breaking Change(s) Detected",
				"major release",
				"`override-breaking-change`",
				"- Breaking change 1 - [reference](doc1)\n- Breaking change 2 - [reference](doc2)\n",
			},
			notExpectedStrings: []string{
				"generated some diffs",
				"## Errors",
				"## Missing test report",
			},
		},
		"missing tests are displayed": {
			data: diffCommentData{
				MissingTests: map[string]*MissingTestInfo{
					"resource": {
						Tests:         []string{"test-a", "test-b"},
						SuggestedTest: "x",
					},
				},
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

func TestFileToResource(t *testing.T) {
	cases := map[string]struct {
		path string
		want string
	}{
		// Resource go files
		"files outside services directory are not resources": {
			path: "/google-beta/tpgiamresource/resource_iam_binding.go",
			want: "",
		},
		"non-go files in service directories are not resources": {
			path: "/google-beta/services/firebaserules/resource_firebaserules_release.html.markdown",
			want: "",
		},
		"resource file": {
			path: "/google-beta/services/firebaserules/resource_firebaserules_release.go",
			want: "google_firebaserules_release",
		},
		"resource iam file": {
			path: "/google/services/kms/iam_kms_crypto_key.go",
			want: "google_kms_crypto_key",
		},
		"resource generated test file": {
			path: "/google-beta/services/containeraws/resource_container_aws_node_pool_generated_test.go",
			want: "google_container_aws_node_pool",
		},
		"resource handwritten test file": {
			path: "/google-beta/services/oslogin/resource_os_login_ssh_public_key_test.go",
			want: "google_os_login_ssh_public_key",
		},
		"resource internal_test file": {
			path: "/google/services/redis/resource_redis_instance_internal_test.go",
			want: "google_redis_instance",
		},
		"resource sweeper file": {
			path: "/google-beta/services/sql/resource_sql_source_representation_instance_sweeper.go",
			want: "google_sql_source_representation_instance",
		},
		"resource iam handwritten test file": {
			path: "/google-beta/services/bigtable/resource_bigtable_instance_iam_test.go",
			want: "google_bigtable_instance",
		},
		"resource iam generated test file": {
			path: "/google-beta/services/privateca/iam_privateca_ca_pool_generated_test.go",
			want: "google_privateca_ca_pool",
		},
		"resource ignore google_ prefix": {
			path: "/google-beta/services/resourcemanager/resource_google_project_sweeper.go",
			want: "google_project",
		},
		"resource starting with iam_": {
			path: "/google-beta/services/iam2/resource_iam_access_boundary_policy.go",
			want: "google_iam_access_boundary_policy",
		},
		"resource file without starting slash": {
			path: "google-beta/services/firebaserules/resource_firebaserules_release.go",
			want: "google_firebaserules_release",
		},

		// Datasource files
		"datasource file": {
			path: "/google/services/dns/data_source_dns_keys.go",
			want: "google_dns_keys",
		},
		"datasource handwritten test file": {
			path: "/google-beta/services/monitoring/data_source_monitoring_service_test.go",
			want: "google_monitoring_service",
		},
		// Future-proofing
		"datasource generated test file": {
			path: "/google-beta/services/alloydb/data_source_alloydb_locations_generated_test.go",
			want: "google_alloydb_locations",
		},
		"datasource internal_test file": {
			path: "/google/services/storage/data_source_storage_object_signed_url_internal_test.go",
			want: "google_storage_object_signed_url",
		},
		"datasource ignore google_ prefix": {
			path: "/google-beta/services/certificatemanager/data_source_google_certificate_manager_certificate_map_test.go",
			want: "google_certificate_manager_certificate_map",
		},
		"datasource starting with iam_": {
			path: "/google-beta/services/resourcemanager/data_source_iam_policy_test.go",
			want: "google_iam_policy",
		},
		"datasource file without starting slash": {
			path: "google/services/dns/data_source_dns_keys.go",
			want: "google_dns_keys",
		},

		// Resource documentation
		"files outside /r or /d directories are not resources": {
			path: "/website/docs/guides/common_issues.html.markdown",
			want: "",
		},
		"non-markdown files are not resources": {
			path: "/website/docs/r/access_context_manager_access_level.go",
			want: "",
		},
		"resource docs": {
			path: "/website/docs/r/firestore_document.html.markdown",
			want: "google_firestore_document",
		},
		"resource docs ignore google_ prefix": {
			path: "/website/docs/r/google_project_service.html.markdown",
			want: "google_project_service",
		},
		"resource docs starting with iam_": {
			path: "/website/docs/r/iam_deny_policy.html.markdown",
			want: "google_iam_deny_policy",
		},
		"resource docs without starting slash": {
			path: "website/docs/d/cloudbuild_trigger.html.markdown",
			want: "google_cloudbuild_trigger",
		},

		// Datasource documentation
		"datasource docs": {
			path: "/website/docs/d/beyondcorp_app_gateway.html.markdown",
			want: "google_beyondcorp_app_gateway",
		},
		"datasource docs ignore google_ prefix": {
			path: "/website/docs/d/google_vertex_ai_index.html.markdown",
			want: "google_vertex_ai_index",
		},
		"datasource docs starting with iam_": {
			path: "/website/docs/d/iam_role.html.markdown",
			want: "google_iam_role",
		},
		"datasource docs without starting slash": {
			path: "website/docs/d/beyondcorp_app_gateway.html.markdown",
			want: "google_beyondcorp_app_gateway",
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			got := fileToResource(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPathChanged(t *testing.T) {
	cases := map[string]struct {
		path         string
		changedFiles []string
		want         bool
	}{
		"no changed files": {
			path:         "path/to/folder/file.go",
			changedFiles: []string{},
			want:         false,
		},
		"path matches exactly": {
			path:         "path/to/folder/file.go",
			changedFiles: []string{"path/to/folder/file.go"},
			want:         true,
		},
		"path matches files in a folder": {
			path:         "path/to/folder/",
			changedFiles: []string{"path/to/folder/file.go"},
			want:         true,
		},
		"path matches partial folder name": {
			path:         "path/to/folder",
			changedFiles: []string{"path/to/folder2/file.go"},
			want:         true,
		},
		"path matches second item in list": {
			path:         "path/to/folder/",
			changedFiles: []string{"path/to/folder2/file.go", "path/to/folder/file.go"},
			want:         true,
		},
		"path doesn't match files in a different folder": {
			path:         "path/to/folder/",
			changedFiles: []string{"path/to/folder2/file.go"},
			want:         false,
		},
		"path doesn't match multiple items": {
			path:         "path/to/folder/",
			changedFiles: []string{"path/to/folder2/file.go", "path/to/folder3"},
			want:         false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			got := pathChanged(tc.path, tc.changedFiles)
			assert.Equal(t, tc.want, got)
		})
	}
}
