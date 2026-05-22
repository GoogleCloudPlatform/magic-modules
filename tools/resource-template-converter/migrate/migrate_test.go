package migrate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateFile_Public(t *testing.T) {
	// Set up a temporary directory structure simulating public magic-modules
	tmpDir, err := ioutil.TempDir("", "mm-public-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create products and templates subdirectories
	productsDir := filepath.Join(tmpDir, "mmv1", "products", "accesscontextmanager")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		t.Fatalf("failed to create products dir: %v", err)
	}

	examplesDir := filepath.Join(tmpDir, "mmv1", "templates", "terraform", "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("failed to create examples dir: %v", err)
	}

	// Create resource YAML file
	yamlPath := filepath.Join(productsDir, "AccessLevel.yaml")
	yamlContent := `---
name: AccessLevel
description: An AccessLevel is a label.
examples:
  - name: access_context_manager_access_level_basic
    primary_resource_id: access-level
    vars:
      access_level_name: chromeos_no_lock
    exclude_test: true
`
	if err := ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	// Create mock example template
	tmplPath := filepath.Join(examplesDir, "access_context_manager_access_level_basic.tf.tmpl")
	tmplContent := `resource "google_access_context_manager_access_level" "access-level" {
  title = "$.Vars.access_level_name"
}`
	if err := ioutil.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write tmpl file: %v", err)
	}

	// Run migration
	err = MigrateFile(yamlPath, "accesscontextmanager")
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify YAML content changes
	updatedYamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read updated yaml: %v", err)
	}
	updatedYaml := string(updatedYamlBytes)

	if !strings.Contains(updatedYaml, "samples:") {
		t.Errorf("expected samples block, got: %s", updatedYaml)
	}
	if strings.Contains(updatedYaml, "examples:") {
		t.Errorf("expected examples block to be removed, got: %s", updatedYaml)
	}
	if !strings.Contains(updatedYaml, "resource_id_vars:") {
		t.Errorf("expected resource_id_vars, got: %s", updatedYaml)
	}
}

func TestMigrateFile_PrivateOverrides(t *testing.T) {
	// Set up a temporary directory structure simulating EAP private overrides
	tmpDir, err := ioutil.TempDir("", "mm-private-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create products and examples subdirectories directly under root
	productsDir := filepath.Join(tmpDir, "products", "accesscontextmanager")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		t.Fatalf("failed to create products dir: %v", err)
	}

	examplesDir := filepath.Join(tmpDir, "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("failed to create examples dir: %v", err)
	}

	// Create resource YAML file
	yamlPath := filepath.Join(productsDir, "ServicePerimeter.yaml")
	yamlContent := `---
name: ServicePerimeter
min_version: private
examples:
  - name: access_context_manager_service_perimeter_weakened_for_testing
    config_path: examples/access_context_manager_service_perimeter_weakened_for_testing.tf.tmpl
    primary_resource_id: service-perimeter
    vars:
      service_perimeter_name: restrict_all
    exclude_test: true
`
	if err := ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	// Create mock example template
	tmplPath := filepath.Join(examplesDir, "access_context_manager_service_perimeter_weakened_for_testing.tf.tmpl")
	tmplContent := `resource "google_access_context_manager_service_perimeter" "service-perimeter" {
  title = "$.Vars.service_perimeter_name"
}`
	if err := ioutil.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write tmpl file: %v", err)
	}

	// Run migration
	err = MigrateFile(yamlPath, "accesscontextmanager")
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify YAML content changes
	updatedYamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read updated yaml: %v", err)
	}
	updatedYaml := string(updatedYamlBytes)

	if !strings.Contains(updatedYaml, "samples:") {
		t.Errorf("expected samples block, got: %s", updatedYaml)
	}
	if strings.Contains(updatedYaml, "examples:") {
		t.Errorf("expected examples block to be removed, got: %s", updatedYaml)
	}
	if !strings.Contains(updatedYaml, "resource_id_vars:") {
		t.Errorf("expected resource_id_vars, got: %s", updatedYaml)
	}
}

func TestMigrateFile_PreserveComments(t *testing.T) {
	// Set up a temporary directory structure
	tmpDir, err := ioutil.TempDir("", "mm-comments-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	productsDir := filepath.Join(tmpDir, "mmv1", "products", "accesscontextmanager")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		t.Fatalf("failed to create products dir: %v", err)
	}

	examplesDir := filepath.Join(tmpDir, "mmv1", "templates", "terraform", "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("failed to create examples dir: %v", err)
	}

	// Create resource YAML file with comments inside examples
	yamlPath := filepath.Join(productsDir, "AccessLevel.yaml")
	yamlContent := `---
name: AccessLevel
examples:
  # This is a comment for the access level basic example mapping
  - name: access_context_manager_access_level_basic
    # This comment is on the primary_resource_id key
    primary_resource_id: access-level
    vars:
      # This comment is inside vars
      access_level_name: chromeos_no_lock
    exclude_test: true
`
	if err := ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	tmplPath := filepath.Join(examplesDir, "access_context_manager_access_level_basic.tf.tmpl")
	if err := ioutil.WriteFile(tmplPath, []byte(`resource "google" "test" {}`), 0644); err != nil {
		t.Fatalf("failed to write tmpl file: %v", err)
	}

	// Run migration
	err = MigrateFile(yamlPath, "accesscontextmanager")
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify YAML content changes and comments preservation
	updatedYamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read updated yaml: %v", err)
	}
	updatedYaml := string(updatedYamlBytes)

	if !strings.Contains(updatedYaml, "# This is a comment for the access level basic example mapping") {
		t.Errorf("expected mapping comment to be preserved, got:\n%s", updatedYaml)
	}
	if !strings.Contains(updatedYaml, "# This comment is on the primary_resource_id key") {
		t.Errorf("expected key comment to be preserved, got:\n%s", updatedYaml)
	}
	if !strings.Contains(updatedYaml, "# This comment is inside vars") {
		t.Errorf("expected vars comment to be preserved, got:\n%s", updatedYaml)
	}
}
