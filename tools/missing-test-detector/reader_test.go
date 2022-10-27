package main

import (
	"testing"
)

func TestReadTestFile(t *testing.T) {
	endpointTests, err := readTestFile("testdata/resource_vertex_ai_endpoint_test.go")
	if err != nil {
		t.Errorf("error reading endpoint test file: %v", err)
	}
	if len(endpointTests) != 1 {
		t.Errorf("unexpected number of endpointTests: %d, expected 1", len(endpointTests))
	}
	if len(endpointTests[0].Steps) != 2 {
		t.Errorf("unexpected number of test steps: %d, expected 2", len(endpointTests[0].Steps))
	}
	if endpoints, ok := endpointTests[0].Steps[0]["google_vertex_ai_endpoint"]; !ok {
		t.Errorf("did not find an endpoint resource in %v", endpointTests[0].Steps[0])
	} else if endpointsMap, ok := endpoints.(map[string]any); !ok {
		t.Errorf("did not find a map of endpoint resources, found %v", endpoints)
	} else if endpoint, ok := endpointsMap["endpoint"]; !ok {
		t.Errorf("did not find an endpoint in %v", endpointsMap)
	} else if endpointConfig, ok := endpoint.(map[string]any); !ok {
		t.Errorf("endpoint config was not a map, was %v", endpoint)
	} else if len(endpointConfig) != 8 {
		t.Errorf("found wrong number of fields in endpoint config: %d, expected 8", len(endpointConfig))
	}
	instanceTests, err := readTestFile("testdata/resource_sql_database_instance_test.go")
	if err != nil {
		t.Errorf("error reading sql database instance test file: %v", err)
	}
	if len(instanceTests) != 37 {
		t.Errorf("unexpected number of instanceTests: %d, expected 37", len(instanceTests))
	}
	for _, instanceTest := range instanceTests {
		for _, step := range instanceTest.Steps {
			if instances, ok := step["google_sql_database_instance"]; !ok {
				t.Errorf("did not find a database instance resource in %v", step)
			} else if instancesMap, ok := instances.(map[string]any); !ok {
				t.Errorf("did not find a map of database instance resources, found %v", instances)
			} else {
				for name, instance := range instancesMap {
					if _, ok := instance.(map[string]any); !ok {
						t.Errorf("database instance %s config was not a map, was %v", name, instance)
					}
				}
			}
		}
	}
}
