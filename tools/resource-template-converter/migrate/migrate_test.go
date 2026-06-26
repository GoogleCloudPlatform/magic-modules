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
	err = MigrateFile(yamlPath, "accesscontextmanager", false, false)
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
	err = MigrateFile(yamlPath, "accesscontextmanager", false, false)
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
	err = MigrateFile(yamlPath, "accesscontextmanager", false, false)
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

func TestMigrateFile_DiscardPrimaryResourceName(t *testing.T) {
	// Set up a temporary directory structure
	tmpDir, err := ioutil.TempDir("", "mm-discard-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	productsDir := filepath.Join(tmpDir, "mmv1", "products", "datacatalog")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		t.Fatalf("failed to create products dir: %v", err)
	}

	examplesDir := filepath.Join(tmpDir, "mmv1", "templates", "terraform", "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("failed to create examples dir: %v", err)
	}

	// Create resource YAML file with primary_resource_name inside examples
	yamlPath := filepath.Join(productsDir, "PolicyTag.yaml")
	yamlContent := `---
name: PolicyTag
examples:
  - name: data_catalog_taxonomies_policy_tag_basic
    primary_resource_id: basic_policy_tag
    primary_resource_name: fmt.Sprintf("tf_test_my_policy_tag%s", context["random_suffix"])
    vars:
      taxonomy_display_name: taxonomy_display_name
`
	if err := ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	tmplPath := filepath.Join(examplesDir, "data_catalog_taxonomies_policy_tag_basic.tf.tmpl")
	if err := ioutil.WriteFile(tmplPath, []byte(`resource "google" "test" {}`), 0644); err != nil {
		t.Fatalf("failed to write tmpl file: %v", err)
	}

	// Run migration
	err = MigrateFile(yamlPath, "datacatalog", false, false)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify YAML content and that primary_resource_name was discarded
	updatedYamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read updated yaml: %v", err)
	}
	updatedYaml := string(updatedYamlBytes)

	if strings.Contains(updatedYaml, "primary_resource_name") {
		t.Errorf("expected primary_resource_name to be discarded, but it was present in the migrated YAML:\n%s", updatedYaml)
	}
}

func TestMigrateFile_DiscardUnrecognizedFields(t *testing.T) {
	// Set up a temporary directory structure
	tmpDir, err := ioutil.TempDir("", "mm-unrecognized-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	productsDir := filepath.Join(tmpDir, "mmv1", "products", "artifactregistry")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		t.Fatalf("failed to create products dir: %v", err)
	}

	examplesDir := filepath.Join(tmpDir, "mmv1", "templates", "terraform", "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("failed to create examples dir: %v", err)
	}

	// Create resource YAML file with exclude_from_docs inside examples
	yamlPath := filepath.Join(productsDir, "Rule.yaml")
	yamlContent := `---
name: Rule
examples:
  - name: artifact_registry_rule_full
    primary_resource_id: my-rule
    exclude_from_docs: true
    vars:
      repository_id: my-repository
`
	if err := ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	tmplPath := filepath.Join(examplesDir, "artifact_registry_rule_full.tf.tmpl")
	if err := ioutil.WriteFile(tmplPath, []byte(`resource "google" "test" {}`), 0644); err != nil {
		t.Fatalf("failed to write tmpl file: %v", err)
	}

	// Run migration
	err = MigrateFile(yamlPath, "artifactregistry", false, false)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify YAML content and that exclude_from_docs was discarded
	updatedYamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read updated yaml: %v", err)
	}
	updatedYaml := string(updatedYamlBytes)

	if strings.Contains(updatedYaml, "exclude_from_docs") {
		t.Errorf("expected exclude_from_docs to be discarded, but it was present in the migrated YAML:\n%s", updatedYaml)
	}
}

func TestMigrateFile_OnlyMigration(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "mm-only-migration-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	productsDir := filepath.Join(tmpDir, "mmv1", "products", "accesscontextmanager")
	os.MkdirAll(productsDir, 0755)

	examplesDir := filepath.Join(tmpDir, "mmv1", "templates", "terraform", "examples")
	os.MkdirAll(examplesDir, 0755)

	// Create resource YAML file with unordered keys and quotes
	yamlPath := filepath.Join(productsDir, "AccessLevel.yaml")
	yamlContent := `---
# Access level yaml
name: "AccessLevel"
description: "An AccessLevel is a label."
examples:
  - name: access_context_manager_access_level_basic
    primary_resource_id: access-level
    vars:
      access_level_name: chromeos_no_lock
    exclude_test: true
`
	ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644)

	tmplPath := filepath.Join(examplesDir, "access_context_manager_access_level_basic.tf.tmpl")
	ioutil.WriteFile(tmplPath, []byte(`resource "google" "test" {}`), 0644)

	// Run migration only
	err = MigrateFile(yamlPath, "accesscontextmanager", true, false)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	updatedYamlBytes, _ := ioutil.ReadFile(yamlPath)
	updatedYaml := string(updatedYamlBytes)

	// Migration should happen:
	if !strings.Contains(updatedYaml, "samples:") {
		t.Errorf("expected samples block, got: %s", updatedYaml)
	}
	// Formatting should NOT happen (quotes should remain):
	if !strings.Contains(updatedYaml, `"AccessLevel"`) {
		t.Errorf("expected string quotes to be preserved under only-migration, got: %s", updatedYaml)
	}
}

func TestMigrateFile_OnlyFormat(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "mm-only-format-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	productsDir := filepath.Join(tmpDir, "mmv1", "products", "accesscontextmanager")
	os.MkdirAll(productsDir, 0755)

	yamlPath := filepath.Join(productsDir, "AccessLevel.yaml")
	yamlContent := `---
# Access level yaml
examples:
  - name: access_context_manager_access_level_basic
    primary_resource_id: access-level
name: "AccessLevel"
description: "An AccessLevel is a label."
`
	ioutil.WriteFile(yamlPath, []byte(yamlContent), 0644)

	// Run formatting only
	err = MigrateFile(yamlPath, "accesscontextmanager", false, true)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	updatedYamlBytes, _ := ioutil.ReadFile(yamlPath)
	updatedYaml := string(updatedYamlBytes)

	// Migration should NOT happen:
	if strings.Contains(updatedYaml, "samples:") {
		t.Errorf("expected samples block to NOT be created, got: %s", updatedYaml)
	}
	if !strings.Contains(updatedYaml, "examples:") {
		t.Errorf("expected examples block to remain, got: %s", updatedYaml)
	}
	// Formatting should happen (quotes stripped):
	if strings.Contains(updatedYaml, `"AccessLevel"`) {
		t.Errorf("expected string quotes to be stripped under only-format, got: %s", updatedYaml)
	}
}
