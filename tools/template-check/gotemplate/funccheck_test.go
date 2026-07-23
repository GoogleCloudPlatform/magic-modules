package gotemplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go.tmpl")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	return path
}

func TestFuncCheck_ValidFunctions(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"mmv1_custom", `{{camelize .Name}} {{underscore .Name}} {{title .Name}}`},
		{"go_builtins", `{{len .Items}} {{and .A .B}} {{index .Map "key"}}`},
		{"keywords", `{{if .A}} {{else}} {{end}} {{range .Items}} {{with .X}}`},
		{"dot_access", `{{.Resource.Name}} {{.Name}}`},
		{"variables", `{{$name := .Name}} {{$name}}`},
		{"string_literals", `{{print "hello"}} {{printf "%s" .Name}}`},
		{"complex_mmv1", `{{plural .Name}} {{contains .A .B}} {{join .List ","}}`},
		{"provider_funcs", `{{TemplatePath "compute"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestFile(t, tt.content)
			results, err := CheckInvalidFuncsForFile(path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results) != 0 {
				t.Errorf("expected no errors, got: %v", results)
			}
		})
	}
}

func TestFuncCheck_InvalidFunctions(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{"typo", `{{camelCase .Name}}`, "camelCase"},
		{"missing_mmv1", `{{toLower .Name}}`, "toLower"},
		{"constant_as_func", `{{BigQueryBasePath}}`, "BigQueryBasePath"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestFile(t, tt.content)
			results, _ := CheckInvalidFuncsForFile(path)
			if len(results) == 0 || !strings.Contains(results[0], tt.expected) {
				t.Errorf("expected error for %s, got: %v", tt.expected, results)
			}
		})
	}
}
