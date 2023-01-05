package main

import (
	"strings"
)

type MissingTestInfo struct {
	UntestedFields []string
	Tests          []string
}

type FieldSet map[string]struct{}

// Detect missing tests for the given resource changes map in the given slice of tests.
// Return a map of resource names to missing test info about that resource.
func detectMissingTests(changedFields map[string]FieldCoverage, allTests []*Test) map[string]*MissingTestInfo {
	resourceNamesToTests := make(map[string][]string)
	for _, test := range allTests {
		for _, step := range test.Steps {
			for resourceName, resourceMap := range step {
				if changedResourceFields, ok := changedFields[resourceName]; ok {
					// This resource type has changed fields.
					resourceNamesToTests[resourceName] = append(resourceNamesToTests[resourceName], test.Name)
					for _, resourceConfig := range resourceMap {
						markCoverage(changedResourceFields, resourceConfig)
					}
				}
			}
		}
	}
	missingTests := make(map[string]*MissingTestInfo)
	for resourceName, fieldCoverage := range changedFields {
		untested := untestedFields(fieldCoverage, nil)
		if len(untested) > 0 {
			missingTests[resourceName] = &MissingTestInfo{
				UntestedFields: untestedFields(fieldCoverage, nil),
				Tests:          resourceNamesToTests[resourceName],
			}
		}
	}
	return missingTests
}

func markCoverage(fieldCoverage FieldCoverage, config Resource) {
	for fieldName, fieldValue := range config {
		if coverage, ok := fieldCoverage[fieldName]; ok {
			if covered, ok := coverage.(bool); ok {
				if !covered {
					fieldCoverage[fieldName] = true
				}
			} else if objectCoverage, ok := coverage.(FieldCoverage); ok {
				if fieldValueConfig, ok := fieldValue.(Resource); ok {
					markCoverage(objectCoverage, fieldValueConfig)
				}
			}
		}
	}
}

func untestedFields(fieldCoverage FieldCoverage, path []string) []string {
	fields := make([]string, 0)
	for fieldName, coverage := range fieldCoverage {
		if covered, ok := coverage.(bool); ok {
			if !covered {
				fields = append(fields, strings.Join(append(path, fieldName), "."))
			}
		} else if objectCoverage, ok := coverage.(FieldCoverage); ok {
			fields = append(fields, untestedFields(objectCoverage, append(path, fieldName))...)
		}
	}
	return fields
}
