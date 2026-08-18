// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestGenerateFile(t *testing.T) {
	mockFS := fstest.MapFS{
		"templates/test.go.tmpl": &fstest.MapFile{
			Data: []byte(`package test
{{if .Content}}
// Some content
{{.Content}}
{{end}}`),
		},
		"templates/empty.go.tmpl": &fstest.MapFile{
			Data: []byte(`{{if .Content}}package test
{{.Content}}{{end}}`),
		},
		"templates/whitespace.go.tmpl": &fstest.MapFile{
			Data: []byte(`
   
{{if .Content}}
{{.Content}}
{{end}}
   
`),
		},
	}

	tempDir := t.TempDir()

	td := NewTemplateData(tempDir, "ga", mockFS)

	tests := []struct {
		name         string
		filePath     string
		templatePath string
		input        any
		goFormat     bool
		templates    []string
		wantWrite    bool
		wantContent  string
	}{
		{
			name:         "standard template with content",
			filePath:     filepath.Join(tempDir, "standard.go"),
			templatePath: "templates/test.go.tmpl",
			input:        map[string]any{"Content": "var x = 1"},
			goFormat:     true,
			templates:    []string{"templates/test.go.tmpl"},
			wantWrite:    true,
			wantContent:  "package test\n\n// Some content\nvar x = 1\n", // formatted
		},
		{
			name:         "empty template output",
			filePath:     filepath.Join(tempDir, "empty.go"),
			templatePath: "templates/empty.go.tmpl",
			input:        map[string]any{"Content": ""},
			goFormat:     true,
			templates:    []string{"templates/empty.go.tmpl"},
			wantWrite:    false,
		},
		{
			name:         "whitespace-only template output",
			filePath:     filepath.Join(tempDir, "whitespace.go"),
			templatePath: "templates/whitespace.go.tmpl",
			input:        map[string]any{"Content": ""},
			goFormat:     true,
			templates:    []string{"templates/whitespace.go.tmpl"},
			wantWrite:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			td.GenerateFile(tc.filePath, tc.templatePath, tc.input, tc.goFormat, tc.templates...)

			_, err := os.Stat(tc.filePath)
			exists := !os.IsNotExist(err)

			if tc.wantWrite != exists {
				t.Fatalf("expected write: %t, got: %t", tc.wantWrite, exists)
			}

			if tc.wantWrite {
				content, err := os.ReadFile(tc.filePath)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if string(content) != tc.wantContent {
					t.Errorf("expected content:\n%q\ngot:\n%q", tc.wantContent, string(content))
				}
			}
		})
	}
}

func TestGenerateQueryTestFileImportsFirstSampleDependenciesOnly(t *testing.T) {
	productVersion := &product.Version{Name: "ga"}
	importPath := "github.com/hashicorp/terraform-provider-google/google"
	res := api.Resource{
		Name:            "Foo",
		ImportPath:      importPath,
		ProductMetadata: &api.Product{Name: "Example", Version: productVersion, Versions: []*product.Version{productVersion}},
		Samples: []*resource.Sample{
			{
				Name:              "example_foo",
				PrimaryResourceId: "example_foo",
				Steps: []*resource.Step{{
					Name: "example_foo",
					TestContextVars: map[string]string{
						"network_name": `servicenetworking.BootstrapSharedServiceNetworkingConnection(t, "test-network")`,
					},
				}},
			},
			{
				Name:              "example_foo_iam",
				PrimaryResourceId: "example_foo_iam",
				BootstrapIam: []resource.IamMember{{
					Member: "serviceAccount:test@example.com",
					Role:   "roles/viewer",
				}},
				Steps: []*resource.Step{{
					Name: "example_foo_iam",
				}},
			},
		},
	}
	filePath := filepath.Join(t.TempDir(), "list_google_example_foo_generated_test.go")
	NewTemplateData(t.TempDir(), "ga", os.DirFS("..")).GenerateQueryTestFile(filePath, res)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("query file should be generated: %v", err)
	}
	generated := string(content)
	if !strings.Contains(generated, importPath+"/services/servicenetworking") {
		t.Fatalf("query test missing first-sample TestContextVars import:\n%s", generated)
	}
	if strings.Contains(generated, importPath+"/services/resourcemanager") {
		t.Fatalf("query test imported later-sample BootstrapIam dependency:\n%s", generated)
	}
}

func TestGenerateQueryTestFileUsesResourcemanagerForFirstSampleBootstrapIam(t *testing.T) {
	productVersion := &product.Version{Name: "ga"}
	importPath := "github.com/hashicorp/terraform-provider-google/google"
	res := api.Resource{
		Name:            "Foo",
		ImportPath:      importPath,
		ProductMetadata: &api.Product{Name: "Example", Version: productVersion, Versions: []*product.Version{productVersion}},
		Samples: []*resource.Sample{{
			Name:              "example_foo",
			PrimaryResourceId: "example_foo",
			BootstrapIam: []resource.IamMember{{
				Member: "serviceAccount:test@example.com",
				Role:   "roles/viewer",
			}},
			Steps: []*resource.Step{{
				Name: "example_foo",
			}},
		}},
	}
	filePath := filepath.Join(t.TempDir(), "list_google_example_foo_generated_test.go")
	NewTemplateData(t.TempDir(), "ga", os.DirFS("..")).GenerateQueryTestFile(filePath, res)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("query file should be generated: %v", err)
	}
	generated := string(content)
	if !strings.Contains(generated, importPath+"/services/resourcemanager") {
		t.Fatalf("query test missing resourcemanager import for first-sample BootstrapIam:\n%s", generated)
	}
	if !strings.Contains(generated, "resourcemanager.BootstrapIamMembers") {
		t.Fatalf("query test should call resourcemanager.BootstrapIamMembers:\n%s", generated)
	}
	if strings.Contains(generated, "acctest.BootstrapIamMembers") {
		t.Fatalf("query test still calls acctest.BootstrapIamMembers:\n%s", generated)
	}
}

func TestGenerateListResourceStripsSelfLinkFromIdentityName(t *testing.T) {
	productVersion := &product.Version{Name: "ga", BaseUrl: "https://example.googleapis.com/v1/"}
	res := api.Resource{
		Name:                 "Foo",
		GenerateListResource: true,
		ImportPath:           "github.com/hashicorp/terraform-provider-google/google",
		ProductMetadata:      &api.Product{Name: "Example", Versions: []*product.Version{productVersion}, Version: productVersion},
		CollectionUrlKey:     "foos",
		BaseUrl:              "projects/{{project}}/locations/{{location}}/foos",
		IdFormat:             "projects/{{project}}/locations/{{location}}/foos/{{name}}",
		ImportFormat:         []string{"projects/{{project}}/locations/{{location}}/foos/{{name}}"},
		Parameters: []*api.Type{
			{Name: "location", Type: "String", UrlParamOnly: true, Required: true, ApiName: "location"},
			{Name: "name", Type: "String", UrlParamOnly: true, Required: true, ApiName: "name"},
		},
		Properties: []*api.Type{
			{Name: "description", Type: "String", ApiName: "description"},
		},
	}
	filePath := filepath.Join(t.TempDir(), "list_google_example_foo.go")
	NewTemplateData(t.TempDir(), "ga", os.DirFS("..")).GenerateFile(filePath, "templates/terraform/list_resource.go.tmpl", res, true,
		"templates/terraform/list_resource.go.tmpl",
		"templates/terraform/list_resource_method.go.tmpl",
	)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read generated list resource: %v", err)
	}
	generated := string(content)
	if !strings.Contains(generated, "tpgresource.GetResourceNameFromSelfLink(s)") {
		t.Fatalf("list flattener did not strip self-link from identity name:\n%s", generated)
	}
}
