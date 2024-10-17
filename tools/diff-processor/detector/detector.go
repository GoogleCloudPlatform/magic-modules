package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/documentparser"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

type MissingTestInfo struct {
	UntestedFields []string
	SuggestedTest  string
	Tests          []string
}

type FieldSet map[string]struct{}

// ResourceChanges is a nested map with field names as keys and Field objects
// as bottom-level values.
// Fields are assumed not to be covered until detected in a test.
type ResourceChanges map[string]*Field

type Field struct {
	// Added is true when the field is newly added between oldProvider and newProvider.
	Added bool
	// Changed is true when the field type has changed between oldProvider and newProvider.
	Changed bool
	// Tested is true when a test has been found that includes the field.
	Tested bool
}

type MissingDocField struct {
	Field    string
	Section  string
	FilePath string
}

// Detect missing tests for the given resource changes map in the given slice of tests.
// Return a map of resource names to missing test info about that resource.
func DetectMissingTests(schemaDiff diff.SchemaDiff, allTests []*reader.Test) (map[string]*MissingTestInfo, error) {
	changedFields := getChangedFieldsFromSchemaDiff(schemaDiff)
	return getMissingTestsForChanges(changedFields, allTests)
}

// Convert SchemaDiff object to map of ResourceChanges objects.
// Also remove parent fields and output-only fields.
func getChangedFieldsFromSchemaDiff(schemaDiff diff.SchemaDiff) map[string]ResourceChanges {
	changedFields := make(map[string]ResourceChanges)
	for resource, resourceDiff := range schemaDiff {
		resourceChanges := make(ResourceChanges)
		for field, fieldDiff := range resourceDiff.Fields {
			if field == "project" {
				// Skip the project field.
				continue
			}
			if strings.Contains(resource, "iam") && field == "condition" {
				// Skip the condition field of iam resources because some iam resources do not support it.
				continue
			}
			if fieldDiff.New == nil {
				// Skip deleted fields.
				continue
			}
			if fieldDiff.New.Computed && !fieldDiff.New.Optional {
				// Skip output-only fields.
				continue
			}
			if _, ok := fieldDiff.New.Elem.(*schema.Resource); ok {
				// Skip parent fields.
				continue
			}
			if fieldDiff.Old == nil {
				resourceChanges[field] = &Field{Added: true}
			} else {
				resourceChanges[field] = &Field{Changed: true}
			}
		}
		if len(resourceChanges) > 0 {
			changedFields[resource] = resourceChanges
		}
	}
	return changedFields
}

func getMissingTestsForChanges(changedFields map[string]ResourceChanges, allTests []*reader.Test) (map[string]*MissingTestInfo, error) {
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
		untested := untestedFields(fieldCoverage)
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
	for fieldName := range config {
		if field, ok := fieldCoverage[fieldName]; ok {
			field.Tested = true
		}
	}
	return nil
}

func untestedFields(fieldCoverage ResourceChanges) []string {
	fields := make([]string, 0)
	for key, field := range fieldCoverage {
		if !field.Tested {
			fields = append(fields, key)
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

func DetectMissingDocs(schemaDiff diff.SchemaDiff, repoPath string) (map[string][]MissingDocField, error) {
	missingDocFields := make(map[string][]MissingDocField)
	for resource, resourceDiff := range schemaDiff {
		docFilePath := resourceToDocFile(resource, repoPath)
		var argumentsInDoc, attributesInDoc map[string]bool

		content, err := os.ReadFile(docFilePath)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read resource doc %s: %w", docFilePath, err)
		}
		if err == nil {
			parser := documentparser.NewParser()
			err = parser.Parse(content)
			if err != nil {
				return nil, fmt.Errorf("failed to pass document %s: %w", docFilePath, err)
			}
			argumentsInDoc = listToMap(parser.Arguments())
			attributesInDoc = listToMap(parser.Attributes())
		}

		for field, fieldDiff := range resourceDiff.Fields {
			if !isNewField(fieldDiff) {
				continue
			}
			if isAttribute(fieldDiff) {
				if !attributesInDoc[field] {
					missingDocFields[resource] = append(missingDocFields[resource], MissingDocField{Field: field, Section: "Attributes Reference", FilePath: strings.ReplaceAll(docFilePath, repoPath, "")})
				}
			} else if !argumentsInDoc[field] {
				missingDocFields[resource] = append(missingDocFields[resource], MissingDocField{Field: field, Section: "Arguments Reference", FilePath: strings.ReplaceAll(docFilePath, repoPath, "")})
			}
		}
	}
	return missingDocFields, nil
}

func isAttribute(fieldDiff diff.FieldDiff) bool {
	// for compute_instance, some attributes are not only on top level
	return fieldDiff.New.Computed && !fieldDiff.New.Optional
}

func isNewField(fieldDiff diff.FieldDiff) bool {
	return fieldDiff.Old == nil && fieldDiff.New != nil
}

func resourceToDocFile(resource string, repoPath string) string {
	fileBaseName := strings.TrimPrefix(resource, "google_") + ".html.markdown"
	return filepath.Join(repoPath, "website", "docs", "r", fileBaseName)
}

func listToMap(items []string) map[string]bool {
	m := make(map[string]bool)
	for _, item := range items {
		m[item] = true
	}
	return m
}
