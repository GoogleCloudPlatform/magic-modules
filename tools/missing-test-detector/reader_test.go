package main

import (
	"os"
	"testing"
)

// This test only ensures there isn't a panic reading tests in the provider.
func TestReadAllTests(t *testing.T) {
	if providerDir := os.Getenv("PROVIDER_DIR"); providerDir != "" {
		readAllTests(providerDir)
	}
}

func TestReadTestFile(t *testing.T) {
	coveredResourceTests, err := readTestFile("testdata/covered_resource_test.go")
	if err != nil {
		t.Fatalf("error reading covered resource test file: %v", err)
	}
	if len(coveredResourceTests) != 1 {
		t.Fatalf("unexpected number of coveredResourceTests: %d, expected 1", len(coveredResourceTests))
	}
	if len(coveredResourceTests[0].Steps) != 2 {
		t.Fatalf("unexpected number of test steps: %d, expected 1", len(coveredResourceTests[0].Steps))
	}
	if coveredResources, ok := coveredResourceTests[0].Steps[0]["covered_resource"]; !ok {
		t.Errorf("did not find covered_resource in %v", coveredResourceTests[0].Steps[0])
	} else if coveredResource, ok := coveredResources["resource"]; !ok {
		t.Errorf("did not find a covered resource in %v", coveredResources)
	} else if len(coveredResource) != 2 {
		t.Errorf("found wrong number of fields in covered resource config: %d, expected 2", len(coveredResource))
	}
	configVariableTests, err := readTestFile("testdata/config_variable_test.go")
	if err != nil {
		t.Fatalf("error reading config variable test file: %v", err)
	}
	if len(configVariableTests) != 1 {
		t.Fatalf("unexpected number of instanceTests: %d, expected 1", len(configVariableTests))
	}

}
