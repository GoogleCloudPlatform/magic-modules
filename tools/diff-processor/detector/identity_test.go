package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMissingIdentityCoverage(t *testing.T) {
	// Create a temp directory structure simulating services/
	tmpDir := t.TempDir()
	serviceDir := filepath.Join(tmpDir, "compute")
	os.MkdirAll(serviceDir, 0755)

	// Resource WITH identity, WITH full CRUD coverage
	goodResource := `package compute

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		Identity: &schema.ResourceIdentity{},
	}
}

func resourceComputeInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func resourceComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func resourceComputeInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}
`
	os.WriteFile(filepath.Join(serviceDir, "resource_google_compute_instance.go"), []byte(goodResource), 0644)

	// Test file with import identity test
	goodTest := `package compute

func TestAccComputeInstance_importBlockWithResourceIdentity(t *testing.T) {}
`
	os.WriteFile(filepath.Join(serviceDir, "resource_google_compute_instance_test.go"), []byte(goodTest), 0644)

	// Resource WITH identity, MISSING Create and Update coverage
	badResource := `package compute

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceComputeBadDisk() *schema.Resource {
	return &schema.Resource{
		Identity: &schema.ResourceIdentity{},
	}
}

func resourceComputeBadDiskCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceComputeBadDiskRead(d *schema.ResourceData, meta interface{}) error {
	tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{})
	return nil
}

func resourceComputeBadDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
`
	os.WriteFile(filepath.Join(serviceDir, "resource_google_compute_bad_disk.go"), []byte(badResource), 0644)

	// No test file for bad_disk -> missing import test

	// Run detector scoped to both resources
	changedResources := []string{"google_compute_instance", "google_compute_bad_disk"}
	results, err := DetectMissingIdentityCoverage(tmpDir, changedResources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Good resource should NOT appear in results
	if _, ok := results["google_compute_instance"]; ok {
		t.Error("expected google_compute_instance to pass, but it was flagged")
	}

	// Bad resource SHOULD appear
	bad, ok := results["google_compute_bad_disk"]
	if !ok {
		t.Fatal("expected google_compute_bad_disk to be flagged, but it was not")
	}

	// Should flag Create and Update as missing
	expectedMissing := map[string]bool{"Create": true, "Update": true}
	for _, crud := range bad.MissingCRUD {
		if !expectedMissing[crud] {
			t.Errorf("unexpected missing CRUD: %s", crud)
		}
		delete(expectedMissing, crud)
	}
	for crud := range expectedMissing {
		t.Errorf("expected %s to be flagged as missing but it was not", crud)
	}

	// Should flag missing import test
	if !bad.MissingImportTest {
		t.Error("expected MissingImportTest to be true for google_compute_bad_disk")
	}
}

func TestDetectMissingIdentityCoverage_SkipsUnchangedResources(t *testing.T) {
	tmpDir := t.TempDir()
	serviceDir := filepath.Join(tmpDir, "dns")
	os.MkdirAll(serviceDir, 0755)

	// Resource with identity but missing coverage
	resource := `package dns

func ResourceDnsRecordSet() *schema.Resource {
	return &schema.Resource{
		Identity: &schema.ResourceIdentity{},
	}
}

func resourceDnsRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDnsRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
`
	os.WriteFile(filepath.Join(serviceDir, "resource_dns_record_set.go"), []byte(resource), 0644)

	// Pass empty changedResources - should skip everything
	results, err := DetectMissingIdentityCoverage(tmpDir, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected no results when no resources changed, got %d", len(results))
	}
}
