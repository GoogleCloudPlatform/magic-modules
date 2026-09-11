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

func TestFuncCheck_Basic(t *testing.T) {
	// Test that it catches an invalid function
	path := createTestFile(t, `{{BigQueryBasePath}}`)
	results, err := CheckInvalidFuncsForFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || !strings.Contains(results[0], "BigQueryBasePath") {
		t.Errorf("expected error for BigQueryBasePath, got: %v", results)
	}

	// Test that it passes a valid function
	path2 := createTestFile(t, `{{camelize .Name}}`)
	results2, err := CheckInvalidFuncsForFile(path2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results2) != 0 {
		t.Errorf("expected no errors for camelize, got: %v", results2)
	}
}
