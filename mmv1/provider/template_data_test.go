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
	"regexp"
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

func TestGenerateQueryTestFileCapturesOptionalProjectScope(t *testing.T) {
	productVersion := &product.Version{Name: "ga", BaseUrl: "https://example.googleapis.com/v1/"}
	prod := &api.Product{Name: "Example", Versions: []*product.Version{productVersion}, Version: productVersion}
	res := api.Resource{
		Name:                 "Foo",
		GenerateListResource: true,
		ImportPath:           "github.com/hashicorp/terraform-provider-google/google",
		ProductMetadata:      prod,
		TargetVersionName:    "ga",
		BaseUrl:              "projects/{{project}}/locations/{{location}}/foos",
		IdFormat:             "projects/{{project}}/locations/{{location}}/foos/{{name}}",
		ImportFormat:         []string{"projects/{{project}}/locations/{{location}}/foos/{{name}}"},
		Parameters: []*api.Type{
			{Name: "location", Type: "String", UrlParamOnly: true, Required: true},
			{Name: "name", Type: "String", UrlParamOnly: true, Required: true},
		},
		Properties: []*api.Type{
			{Name: "displayName", Type: "String"},
		},
		Samples: []*resource.Sample{{
			Name:              "example_foo",
			PrimaryResourceId: "example_foo",
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
	if !regexp.MustCompile(`"project":\s+"google_example_foo.example_foo"`).MatchString(generated) {
		t.Fatalf("query test did not capture optional project scope:\n%s", generated)
	}
	if !strings.Contains(generated, "project = var.project") {
		t.Fatalf("query test list config missing project = var.project:\n%s", generated)
	}
	if !strings.Contains(generated, `variable "project" { type = string }`) {
		t.Fatalf("query test missing project input variable:\n%s", generated)
	}
}
