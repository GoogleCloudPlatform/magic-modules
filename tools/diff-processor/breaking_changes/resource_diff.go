package breaking_changes

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

// ResourceDiffRule is a rule that operates on an entire ResourceDiff
type ResourceDiffRule struct {
	Identifier string
	Messages   func(resource string, resourceDiff diff.ResourceDiff) []string
}

// ResourceDiffRules is a list of all ResourceDiff rules
var ResourceDiffRules = []ResourceDiffRule{RemovingAField, AddingExactlyOneOf}

var RemovingAField = ResourceDiffRule{
	Identifier: "resource-schema-field-removal-or-rename",
	Messages:   RemovingAFieldMessages,
}

// TODO: Make field removal a FieldDiffRule b/300124253
func RemovingAFieldMessages(resource string, resourceDiff diff.ResourceDiff) []string {
	fieldsRemoved := []string{}
	for field, fieldDiff := range resourceDiff.Fields {
		if fieldDiff.Old != nil && fieldDiff.New == nil {
			fieldsRemoved = append(fieldsRemoved, field)
		}
	}

	tmpl := "Field `%s` within resource `%s` was either removed or renamed"
	var messages []string
	for _, field := range fieldsRemoved {
		messages = append(messages, fmt.Sprintf(tmpl, field, resource))
	}
	return messages
}

var AddingExactlyOneOf = ResourceDiffRule{
	Identifier: "resource-schema-field-addition-of-exactly-one-of",
	Messages:   AddingExactlyOneOfMessages,
}

func AddingExactlyOneOfMessages(resource string, resourceDiff diff.ResourceDiff) []string {
	var messages []string
	newFieldSets := make(map[string]diff.FieldSet) // Set of field sets in new and not in old.
	oldFieldSets := make(map[string]diff.FieldSet) // Set of field sets in old and not in new.
	for key, fieldSet := range resourceDiff.FieldSets.New.ExactlyOneOf {
		if _, ok := resourceDiff.FieldSets.Old.ExactlyOneOf[key]; !ok {
			newFieldSets[key] = fieldSet
		}
	}
	for key, fieldSet := range resourceDiff.FieldSets.Old.ExactlyOneOf {
		if _, ok := resourceDiff.FieldSets.New.ExactlyOneOf[key]; !ok {
			oldFieldSets[key] = fieldSet
		}
	}
	// Find old field sets which are subsets of new field sets.
	for _, newFieldSet := range newFieldSets {
		var addedFields diff.FieldSet
		found := false
		for _, oldFieldSet := range oldFieldSets {
			if oldFieldSet.IsSubsetOf(newFieldSet) {
				addedFields = newFieldSet.Difference(oldFieldSet)
				found = true
				break
			}
		}
		if !found {
			addedFields = newFieldSet
		}
		for field := range addedFields {
			if fieldDiff, ok := resourceDiff.Fields[field]; ok && fieldDiff.Old != nil && !fieldDiff.Old.Required {
				messages = append(messages, fmt.Sprintf("Field `%s` within resource `%s` was added to exactly one of", field, resource))
			}
		}
	}
	return messages
}
