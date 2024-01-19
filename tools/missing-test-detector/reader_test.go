package main

import (
	"os"
	"reflect"
	"testing"
)

// This test only ensures there isn't a panic reading tests in the provider.
func TestReadAllTests(t *testing.T) {
	if servicesDir := os.Getenv("SERVICES_DIR"); servicesDir != "" {
		_, errs := readAllTests(servicesDir)
		for path, err := range errs {
			t.Logf("path: %s, err: %v", path, err)
		}
	} else {
		t.Log("no services directory provided, skipping TestReadAllTests")
	}
}

func TestReadCoveredResourceTestFile(t *testing.T) {
	tests, err := readTestFiles([]string{"testdata/service/covered_resource_test.go"})
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
				"field_six": "true",
			},
		},
		"field_seven": "true",
	}); !reflect.DeepEqual(coveredResource, expectedResource) {
		t.Errorf("found wrong fields in covered resource config: %#v, expected %#v", coveredResource, expectedResource)
	}
}

func TestReadConfigVariableTestFile(t *testing.T) {
	tests, err := readTestFiles([]string{"testdata/service/config_variable_test.go"})
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
	tests, err := readTestFiles([]string{"testdata/service/multiple_resource_test.go"})
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

func TestReadSerialResourceTestFile(t *testing.T) {
	tests, err := readTestFiles([]string{"testdata/service/serial_resource_test.go"})
	if err != nil {
		t.Fatalf("error reading serial resource test file: %v", err)
	}
	if len(tests) != 2 {
		t.Fatalf("unexpected number of tests: %d, expected 2", len(tests))
	}
	if expectedTests := []*Test{
		{
			Name: "testAccSerialResource1",
			Steps: []Step{
				{
					"serial_resource": {
						"resource": {"field_one": "\"value-one\""},
					},
				},
			},
		},
		{
			Name: "testAccSerialResource2",
			Steps: []Step{
				{
					"serial_resource": {
						"resource": {
							"field_two": Resource{
								"field_three": "\"value-two\"",
							},
						},
					},
				},
			},
		},
	}; !reflect.DeepEqual(tests, expectedTests) {
		t.Errorf("found unexpected serialized tests: %v, expected %v", tests, expectedTests)
	}

}

func TestReadCrossFileTests(t *testing.T) {
	tests, err := readTestFiles([]string{"testdata/service/cross_file_1_test.go", "testdata/service/cross_file_2_test.go"})
	if err != nil {
		t.Fatalf("error reading cross file tests: %v", err)
	}

	expectedTests := []*Test{
		{
			Name: "testAccCrossFile1",
			Steps: []Step{
				{
					"serial_resource": {
						"resource": {"field_one": "\"value-one\""},
					},
				},
			},
		},
		{
			Name: "testAccCrossFile2",
			Steps: []Step{
				{
					"serial_resource": {
						"resource": {
							"field_two": Resource{
								"field_three": "\"value-two\"",
							},
						},
					},
				},
			},
		},
	}

	if len(tests) != len(expectedTests) {
		t.Fatalf("unexpected number of tests: %d, expected %d", len(tests), len(expectedTests))
	}

	if !reflect.DeepEqual(tests, expectedTests) {
		t.Errorf("found unexpected cross file tests: %v, expected %v", tests, expectedTests)
	}

}

func TestReadHelperFunctionCall(t *testing.T) {
	tests, err := readTestFiles([]string{"testdata/service/function_call_test.go"})
	if err != nil {
		t.Fatalf("error reading function call test: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	expectedTest := &Test{
		Name: "TestAccFunctionCallResource",
		Steps: []Step{
			Step{
				"helped_resource": Resources{
					"primary": Resource{
						"field_one": "\"value-one\"",
					},
				},
				"helper_resource": Resources{
					"default": Resource{
						"field_one": "\"value-one\"",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(tests[0], expectedTest) {
		t.Errorf("found unexpected tests using helper function: %v, expected %v", tests[0], expectedTest)
	}
}
