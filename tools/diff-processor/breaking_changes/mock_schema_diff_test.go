package breaking_changes

// MockSchemaDiff implements the diff.ResourceDiffInterface interface for testing
type MockSchemaDiff struct {
	isNewResource                  bool
	fieldsInNewStructure           map[string]bool // Maps field names to whether they're in a new structure
	fieldsWithOldComputedAncestors map[string]bool // Maps field names to whether an ancestor was Computed-only in old schema
}

func (sd MockSchemaDiff) IsNewResource() bool {
	return sd.isNewResource
}

func (sd MockSchemaDiff) IsFieldInNewNestedStructure(field string) bool {
	return sd.fieldsInNewStructure[field]
}

func (sd MockSchemaDiff) HasOldComputedOnlyAncestor(field string) bool {
	return sd.fieldsWithOldComputedAncestors[field]
}

// Create mock schema diffs for testing
var (
	// Mock for existing resource (not new, field not in new structure)
	existingResourceSchemaDiff = MockSchemaDiff{
		isNewResource:        false,
		fieldsInNewStructure: make(map[string]bool),
	}

	// Mock for new resource
	newResourceSchemaDiff = MockSchemaDiff{
		isNewResource:        true,
		fieldsInNewStructure: make(map[string]bool),
	}

	// Mock for field in new nested structure
	fieldInNewStructureSchemaDiff = MockSchemaDiff{
		isNewResource:        false,
		fieldsInNewStructure: map[string]bool{"field": true},
	}

	// Mock for field under an already Computed-only ancestor in the old schema
	fieldUnderOldComputedOnlyAncestorSchemaDiff = MockSchemaDiff{
		isNewResource:                  false,
		fieldsWithOldComputedAncestors: map[string]bool{"field": true},
	}
)
