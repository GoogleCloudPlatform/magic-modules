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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"magician/source"

	"github.com/stretchr/testify/assert"
)

func TestExecGenerateComment(t *testing.T) {
	sb := newSandbox(t)

	originalWorkspace := os.Getenv("WORKSPACE")
	os.Setenv("WORKSPACE", sb.Dir)
	defer func() {
		if originalWorkspace == "" {
			os.Unsetenv("WORKSPACE")
		} else {
			os.Setenv("WORKSPACE", originalWorkspace)
		}
	}()

	originalPath := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s", sb.Dir, originalPath))
	defer os.Setenv("PATH", originalPath)

	magicModulesDir := filepath.Join(sb.Dir, "workspace", "magic-modules", ".ci", "magician")
	os.MkdirAll(magicModulesDir, 0755)
	sb.Runner.PushDir(magicModulesDir)

	for _, repo := range []string{"tpg", "tpgb", "tgc", "tfoics"} {
		os.MkdirAll(filepath.Join(sb.Dir, "workspace", repo), 0755)
	}

	gitScript := `#!/bin/bash
		if [[ "$*" == *"diff"* && "$*" == *"--shortstat"* ]]; then
			if [[ "$PWD" == *"/tpg" || "$PWD" == *"/tpgb" ]]; then
				echo " 2 files changed, 40 insertions(+)"
			elif [[ "$PWD" == *"/tgc" ]]; then
				echo " 1 file changed, 10 insertions(+)"
			fi
		elif [[ "$*" == *"clone"* ]]; then
			exit 0
		fi
		`
	os.WriteFile(filepath.Join(sb.Dir, "git"), []byte(gitScript), 0755)

	makeScript := `#!/bin/bash
mkdir -p bin
cat << 'EOF' > bin/diff-processor
#!/bin/bash
if [[ "$*" == *"schema-diff"* ]]; then
	echo '{"AddedResources": ["google_alloydb_instance"]}'
elif [[ "$*" == *"detect-missing-tests"* ]]; then
	echo '{"google_folder_access_approval_settings":{"SuggestedTest":"resource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}","Tests":["a","b","c"]}}'
elif [[ "$*" == *"detect-missing-docs"* ]]; then
	echo '{"Resource":[],"DataSource":[]}'
fi
EOF
chmod 0755 bin/diff-processor
exit 0
`
	os.WriteFile(filepath.Join(sb.Dir, "make"), []byte(makeScript), 0755)

	gh := &mockGithub{
		calledMethods: make(map[string][][]any),
	}
	ctlr := source.NewController(filepath.Join(sb.Dir, "workspace", "go"), "modular-magician", "*******", sb.Runner)

	for _, repo := range []string{
		"terraform-provider-google",
		"terraform-provider-google-beta",
		"terraform-google-conversion",
	} {
		variablePathOld := filepath.Join(sb.Dir, fmt.Sprintf("commitSHA_modular-magician_%s-old.txt", repo))
		variablePath := filepath.Join(sb.Dir, fmt.Sprintf("commitSHA_modular-magician_%s.txt", repo))
		err := sb.Runner.WriteFile(variablePathOld, "1a2a3a4a")
		if err != nil {
			t.Errorf("Error writing file: %s", err)
		}
		err = sb.Runner.WriteFile(variablePath, "1a2a3a4b")
		if err != nil {
			t.Errorf("Error writing file: %s", err)
		}
	}
	execGenerateComment(
		123456,
		"*******",
		"build1",
		"17",
		"project1",
		"sha1",
		gh,
		sb.Runner,
		ctlr,
	)

	for method, expectedCalls := range map[string][][]any{
		"PostBuildStatus": {
			{"123456", "terraform-provider-multiple-resources", "success", "https://console.cloud.google.com/cloud-build/builds;region=global/build1;step=17?project=project1", "sha1"},
			{"123456", "terraform-provider-breaking-change-test", "success", "https://console.cloud.google.com/cloud-build/builds;region=global/build1;step=17?project=project1", "sha1"},
			{"123456", "terraform-provider-missing-service-labels", "success", "https://console.cloud.google.com/cloud-build/builds;region=global/build1;step=17?project=project1", "sha1"},
		},
		"PostComment": {{"123456", "Hi there, I'm the Modular magician. I've detected the following information about your changes for commit sha1:\n\n## Diff report\n\nYour PR generated the following diffs in downstream repositories:\n\n| Repository | Diff Link | Changes |\n| :--- | :--- | :--- |\n| `google` provider | [View Diff](https://github.com/modular-magician/terraform-provider-google/compare/1a2a3a4a..1a2a3a4b) |  2 files changed, 40 insertions(+) |\n| `google-beta` provider | [View Diff](https://github.com/modular-magician/terraform-provider-google-beta/compare/1a2a3a4a..1a2a3a4b) |  2 files changed, 40 insertions(+) |\n| `terraform-google-conversion` | [View Diff](https://github.com/modular-magician/terraform-google-conversion/compare/1a2a3a4a..1a2a3a4b) |  1 file changed, 10 insertions(+) |\n\n\n\n## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n\n\n"}},
		"AddLabels":   {{"123456", []string{"service/alloydb"}}},
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
				Diffs: []Diff{
					{
						Title:        "Repo 1",
						Repo:         "repo-1",
						ShortStat:    "+1 added, -1 removed",
						CommitSHA:    "1a2a3a4b",
						OldCommitSHA: "1a2a3a4a",
					},
					{
						Title:        "Repo 2",
						Repo:         "repo-2",
						ShortStat:    "+2 added, -2 removed",
						CommitSHA:    "1a2a3a4d",
						OldCommitSHA: "1a2a3a4c",
					},
				},
			},
			expectedStrings: []string{
				"## Diff report",
				"generated the following diffs",
				"| Repository | Diff Link | Changes |",
				"| Repo 1 | [View Diff](https://github.com/modular-magician/repo-1/compare/1a2a3a4a..1a2a3a4b) | +1 added, -1 removed |",
				"| Repo 2 | [View Diff](https://github.com/modular-magician/repo-2/compare/1a2a3a4c..1a2a3a4d) | +2 added, -2 removed |",
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
		"multiple resources are displayed": {
			data: diffCommentData{
				MultipleResources: []string{"google_redis_instance", "google_alloydb_cluster"},
			},
			expectedStrings: []string{
				"## Diff report",
				"## Multiple resources added",
				"`override-multiple-resources`",
				"split it into multiple PRs",
				"`google_redis_instance`, `google_alloydb_cluster`.",
			},
			notExpectedStrings: []string{
				"generated some diffs",
				"## Errors",
				"## Missing test report",
				"## Missing doc report",
				"## Breaking Change(s) Detected",
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
		"missing docs are displayed": {
			data: diffCommentData{
				MissingDocs: &MissingDocsSummary{
					Resource: []MissingDocInfo{
						{
							Name:     "resource-a",
							FilePath: "website/docs/r/resource-a.html.markdown",
							Fields:   []string{"field-a", "field-b"},
						},
						{
							Name:     "resource-b",
							FilePath: "website/docs/r/resource-b.html.markdown",
							Fields:   []string{"field-a", "field-b"},
						},
					},
					DataSource: []MissingDocInfo{
						{
							Name:     "resource-a",
							FilePath: "website/docs/d/resource-a.html.markdown",
						},
						{
							Name:     "resource-b",
							FilePath: "website/docs/d/resource-b.html.markdown",
						},
					},
				},
			},
			expectedStrings: []string{
				"## Diff report",
				"## Missing doc report",
			},
		},
		"missing docs should not be displayed": {
			data: diffCommentData{
				MissingDocs: &MissingDocsSummary{
					Resource:   []MissingDocInfo{},
					DataSource: []MissingDocInfo{},
				},
			},
			notExpectedStrings: []string{
				"## Missing doc report",
			},
		},
		"missing docs should not be displayed when MissingDocs is nil": {
			data: diffCommentData{},
			notExpectedStrings: []string{
				"## Missing doc report",
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

func TestMultipleResources(t *testing.T) {
	cases := []struct {
		name      string
		resources []string
		want      []string
	}{
		{
			name: "no resources",
		},
		{
			name:      "single non-iam",
			resources: []string{"google_redis_instance"},
			want:      []string{"google_redis_instance"},
		},
		{
			name:      "multiple non-iam",
			resources: []string{"google_redis_instance", "google_alloydb_cluster"},
			want:      []string{"google_alloydb_cluster", "google_redis_instance"},
		},
		{
			name:      "single iam only",
			resources: []string{"google_redis_instance_iam_member", "google_redis_instance_iam_policy", "google_redis_instance_iam_binding"},
			want:      []string{"google_redis_instance_iam_*"},
		},
		{
			name:      "single iam with parent",
			resources: []string{"google_redis_instance_iam_member", "google_redis_instance_iam_policy", "google_redis_instance_iam_binding", "google_redis_instance"},
			want:      []string{"google_redis_instance"},
		},
		{
			name: "multiple iam",
			resources: []string{
				"google_redis_instance_iam_member",
				"google_redis_instance_iam_policy",
				"google_redis_instance_iam_binding",
				"google_alloydb_cluster_iam_member",
				"google_alloydb_cluster_iam_policy",
				"google_alloydb_cluster_iam_binding",
			},
			want: []string{"google_alloydb_cluster_iam_*", "google_redis_instance_iam_*"},
		},
		{
			name: "multiple iam with parent",
			resources: []string{
				"google_redis_instance_iam_member",
				"google_redis_instance_iam_policy",
				"google_redis_instance_iam_binding",
				"google_alloydb_cluster_iam_member",
				"google_alloydb_cluster_iam_policy",
				"google_alloydb_cluster_iam_binding",
				"google_redis_instance",
			},
			want: []string{"google_alloydb_cluster_iam_*", "google_redis_instance"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := multipleResources(tc.resources)
			assert.Equal(t, tc.want, got)
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

func TestCheckDocumentFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	files := map[string]string{
		"malformed.markdown": `
subcategory: Example Subcategory
---	
`,
		"sample.markdown": `
---
subcategory: Example Subcategory
---	
`,
		"missingsubcategory.markdown": `
---
random: Example Subcategory
---	
`,
	}

	folderPath := filepath.Join(tmpDir, "website", "docs", "r")
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		fullPath := filepath.Join(folderPath, name)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", name, err)
		}
	}

	// write a file in other folders
	if err := os.WriteFile(filepath.Join(tmpDir, "abc.md"), []byte("random"), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", filepath.Join(tmpDir, "abc.md"), err)
	}

	tests := []struct {
		name         string
		changedFiles []string
		wantErr      bool
	}{
		{
			name:         "not in relevant doc folder",
			changedFiles: []string{"abc.md"},
			wantErr:      false,
		},
		{
			name:         "not markdown files",
			changedFiles: []string{"website/docs/r/abc.txt"},
			wantErr:      false,
		},
		{
			name:         "malformed markdown",
			changedFiles: []string{"website/docs/r/malformed.markdown"},
			wantErr:      true,
		},
		{
			name:         "markdown not exist",
			changedFiles: []string{"website/docs/d/sample.markdown"},
			wantErr:      true,
		},
		{
			name:         "correct format",
			changedFiles: []string{"website/docs/r/sample.markdown"},
			wantErr:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := source.Repo{
				Path:         tmpDir,
				ChangedFiles: tc.changedFiles,
			}
			got := checkDocumentFrontmatter(repo)
			if tc.wantErr && len(got) == 0 {
				t.Errorf("checkDocumentFrontmatter() = %v, want error", got)
			}
			if !tc.wantErr && len(got) > 0 {
				t.Errorf("checkDocumentFrontmatter() = %v, want no error", got)
			}
		})
	}
}
