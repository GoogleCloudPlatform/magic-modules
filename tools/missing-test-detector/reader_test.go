package main

import (
	"os"
	"reflect"
	"testing"
)

// This test only ensures there isn't a panic reading tests in the provider.
func TestReadAllTests(t *testing.T) {
	if providerDir := os.Getenv("PROVIDER_DIR"); providerDir != "" {
		readAllTests(providerDir)
	}
}

func TestReadCoveredResourceTestFile(t *testing.T) {
	tests, err := readTestFile("testdata/covered_resource_test.go")
	if err != nil {
		t.Fatalf("error reading covered resource test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if len(tests[0].Steps) != 2 {
		t.Fatalf("unexpected number of test steps: %d, expected 2", len(tests[0].Steps))
	}
	if coveredResources, ok := tests[0].Steps[0]["covered_resource"]; !ok {
		t.Errorf("did not find covered_resource in %v", tests[0].Steps[0])
	} else if coveredResource, ok := coveredResources["resource"]; !ok {
		t.Errorf("did not find a covered resource in %v", coveredResources)
	} else if expectedResource := (Resource{
		"field_one": "\"value-one\"",
		"field_four": Resource{
			"field_five": Resource{
				"field_six": "\"value-three\"",
			},
		},
	}); !reflect.DeepEqual(coveredResource, expectedResource) {
		t.Errorf("found wrong fields in covered resource config: %#v, expected %#v", coveredResource, expectedResource)
	}
}

func TestReadConfigVariableTestFile(t *testing.T) {
	tests, err := readTestFile("testdata/config_variable_test.go")
	if err != nil {
		t.Fatalf("error reading config variable test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if len(tests[0].Steps) != 1 {
		t.Fatalf("unexpected number of test steps: %d, expected 1", len(tests[0].Steps))
	}
	if configVariableResources, ok := tests[0].Steps[0]["config_variable"]; !ok {
		t.Errorf("did not find config_variable in %v", tests[0].Steps[0])
	} else if configVariableResource, ok := configVariableResources["basic"]; !ok {
		t.Errorf("did not find a resource in %v", configVariableResources)
	} else if expectedResource := (Resource{"field_one": "\"value-one\""}); !reflect.DeepEqual(configVariableResource, expectedResource) {
		t.Errorf("found wrong fields in config variable config: %#v, expected %#v", configVariableResource, expectedResource)
	}
}

func TestReadMultipleResourcesTestFile(t *testing.T) {
	tests, err := readTestFile("testdata/multiple_resource_test.go")
	if err != nil {
		t.Fatalf("error reading multiple resources test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if expectedSteps := []Step{
		{
			"resource_one": {
				"instace_two":  {"field_one": "\"value-one\""},
				"instance_one": {"field_one": "\"value-one\""},
			},
			"resource_two": {
				"instace_one": {"field_one": "\"value-one\""},
				"instace_two": {"field_one": "\"value-one\""},
			},
		},
		{
			"resource_one": {
				"instace_two":  {"field_one": "\"value-two\""},
				"instance_one": {"field_one": "\"value-two\""},
			},
			"resource_two": {
				"instace_one": {"field_one": "\"value-two\""},
				"instace_two": {"field_one": "\"value-two\""},
			},
		},
	}; !reflect.DeepEqual(tests[0].Steps, expectedSteps) {
		t.Errorf("found unexpected test steps for multiple resources: %#v, expected %#v", tests[0].Steps, expectedSteps)
	}
}
