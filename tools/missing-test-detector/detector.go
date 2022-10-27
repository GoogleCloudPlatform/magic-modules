package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MissingTestInfo struct {
	UntestedFields []string
	TestCount      int
	StepCount      int
}

// Detect missing tests for the given resource names in the given provider directory.
// Return a map of resource names to missing test info about that resource.
func detectMissingTests(changedFields map[string][]string, providerDir string) (map[string]*MissingTestInfo, error) {
	missingTests := make(map[string]*MissingTestInfo)
	errs := make([]error, 0)
	for resourceName, fields := range changedFields {
		missingTestInfo, err := detectMissingTest(resourceName, providerDir, fields)
		if err != nil {
			errs = append(errs, err)
		}
		missingTests[resourceName] = missingTestInfo
	}
	if len(errs) > 0 {
		return missingTests, fmt.Errorf("errors detecting missing tests: %v", errs)
	}
	return missingTests, nil
}

func detectMissingTest(resourceName, providerDir string, fields []string) (*MissingTestInfo, error) {
	testFileSuffixes := []string{"_test.go", "_generated_test.go"}
	allTests := make([]*Test, 0)
	errs := make([]error, 0)
	for _, testFileSuffix := range testFileSuffixes {
		testFile := filepath.Join(providerDir, strings.Replace(resourceName, "google", "resource", 1)+testFileSuffix)
		tests, err := readTestFile(testFile)
		if err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
		allTests = append(allTests, tests...)
	}
	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("errors reading test files: %v", errs)
	}
	stepCount := 0
	for _, test := range allTests {
		stepCount += len(test.Steps)
	}
	if untestedFields := compareResourceToTests(resourceName, fields, allTests); len(untestedFields) > 0 {
		return &MissingTestInfo{
			UntestedFields: untestedFields,
			TestCount:      len(allTests),
			StepCount:      stepCount,
		}, err
	}
	return nil, err
}

// Return a list of fields as dot-separated paths in the given resource that are not covered by the given tests.
func compareResourceToTests(resourceName string, fields []string, tests []*Test) []string {
	untestedFields := make([]string, 0)
	for _, field := range fields {
		if !fieldInTests(resourceName, strings.Split(field, "."), tests) {
			untestedFields = append(untestedFields, field)
		}
	}
	return untestedFields
}

// Return true if field is present in at least one step of the given tests.
func fieldInTests(resourceName string, path []string, tests []*Test) bool {
	for _, test := range tests {
		for _, step := range test.Steps {
			if resources, ok := step[resourceName].(map[string]any); ok {
				// Resources is a map of all resources of the given type.
				for _, resource := range resources {
					// Loop through all instances of the main tested resource.
					field := resource.(map[string]any)
					present := true
					for _, fieldName := range path {
						fieldValue, ok := field[fieldName]
						if !ok {
							present = false
							break
						}
						field, ok = fieldValue.(map[string]any)
						if !ok {
							break
						}
					}
					if present {
						return true
					}
				}
			}
		}
	}
	return false
}
