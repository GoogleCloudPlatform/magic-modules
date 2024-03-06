package main

import (
	"fmt"
	"sort"
	"strings"

	"missing-test-detector/reader"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type MissingTestInfo struct {
	UntestedFields []string
	SuggestedTest  string
	Tests          []string
}

type FieldSet map[string]struct{}

// Detect missing tests for the given resource changes map in the given slice of tests.
// Return a map of resource names to missing test info about that resource.
func detectMissingTests(changedFields map[string]ResourceChanges, allTests []*reader.Test) (map[string]*MissingTestInfo, error) {
	resourceNamesToTests := make(map[string][]string)
	for _, test := range allTests {
		for _, step := range test.Steps {
			for resourceName, resourceMap := range step {
				if changedResourceFields, ok := changedFields[resourceName]; ok {
					// This resource type has changed fields.
					resourceNamesToTests[resourceName] = append(resourceNamesToTests[resourceName], test.Name)
					for _, resourceConfig := range resourceMap {
						if err := markCoverage(changedResourceFields, resourceConfig); err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	missingTests := make(map[string]*MissingTestInfo)
	for resourceName, fieldCoverage := range changedFields {
		untested := untestedFields(fieldCoverage, nil)
		sort.Strings(untested)
		if len(untested) > 0 {
			missingTests[resourceName] = &MissingTestInfo{
				UntestedFields: untested,
				SuggestedTest:  suggestedTest(resourceName, untested),
				Tests:          resourceNamesToTests[resourceName],
			}
		}
	}
	return missingTests, nil
}

func markCoverage(fieldCoverage ResourceChanges, config reader.Resource) error {
	for fieldName, fieldValue := range config {
		if coverage, ok := fieldCoverage[fieldName]; ok {
			if field, ok := coverage.(*Field); ok {
				field.Tested = true
			} else if objectCoverage, ok := coverage.(ResourceChanges); ok {
				if fieldValueConfig, ok := fieldValue.(reader.Resource); ok {
					if err := markCoverage(objectCoverage, fieldValueConfig); err != nil {
						return fmt.Errorf("error parsing %q: %s", fieldName, err)
					}
				}
			} else {
				return fmt.Errorf("found unexpected %T in field %q", coverage, fieldName)
			}
		}
	}
	return nil
}

func untestedFields(fieldCoverage ResourceChanges, path []string) []string {
	fields := make([]string, 0)
	for fieldName, coverage := range fieldCoverage {
		if field, ok := coverage.(*Field); ok {
			if !field.Tested {
				fields = append(fields, strings.Join(append(path, fieldName), "."))
			}
		} else if objectCoverage, ok := coverage.(ResourceChanges); ok {
			fields = append(fields, untestedFields(objectCoverage, append(path, fieldName))...)
		}
	}
	return fields
}

func suggestedTest(resourceName string, untested []string) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	resourceBlock := rootBody.AppendNewBlock("resource", []string{resourceName, "primary"})
	for _, field := range untested {
		body := resourceBlock.Body()
		path := strings.Split(field, ".")
		for i, step := range path {
			if i < len(path)-1 {
				block := body.FirstMatchingBlock(step, nil)
				if block == nil {
					block = body.AppendNewBlock(step, nil)
				}
				body = block.Body()
			} else {
				body.SetAttributeValue(step, cty.StringVal("VALUE"))
			}
		}
	}
	return strings.ReplaceAll(string(f.Bytes()), `"VALUE"`, "# value needed")
}
