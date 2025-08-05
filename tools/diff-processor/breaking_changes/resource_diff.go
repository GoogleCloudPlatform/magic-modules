package breaking_changes

import (
	"fmt"
	"strings"

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
	messages := []string{}

	for newKey, newSet := range resourceDiff.FieldSets.New.ExactlyOneOf {
		if _, ok := resourceDiff.FieldSets.Old.ExactlyOneOf[newKey]; ok {
			continue // Unchanged EOO.
		}

		// Determine the type of change.
		isSimpleModification := false
		var simpleAddedFields diff.FieldSet

		for _, oldSet := range resourceDiff.FieldSets.Old.ExactlyOneOf {
			if oldSet.IsSubsetOf(newSet) {
				isSimpleModification = true
				simpleAddedFields = newSet.Difference(oldSet)
				break
			}
		}

		if isSimpleModification {
			// Simple modification: only added fields to an existing EOO.
			// Only added *existing* optional fields are breaking.
			for field := range simpleAddedFields {
				if !isNewField(field, resourceDiff) && !isExistingFieldRequired(field, resourceDiff) {
					messages = append(messages, fmt.Sprintf("Field `%s` within resource `%s` was added to exactly one of", field, resource))
				}
			}
		} else if isComplexModification(newSet, resourceDiff) {
			// Complex modification: e.g., add and remove.
			// Any existing, optional field in the new set is breaking. New fields are not.
			for field := range newSet {
				if !isNewField(field, resourceDiff) && !isExistingFieldRequired(field, resourceDiff) {
					messages = append(messages, fmt.Sprintf("Field `%s` within resource `%s` was added to exactly one of", field, resource))
				}
			}
		} else {
			// Brand new EOO.
			// Not breaking if it relaxes a previously required field.
			isRelaxingRequired := false
			for field := range newSet {
				if isExistingFieldRequired(field, resourceDiff) {
					isRelaxingRequired = true
					break
				}
			}
			if isRelaxingRequired {
				continue
			}

			// Not breaking if all fields are in a new optional ancestor.
			isContained := true
			if len(newSet) == 0 {
				isContained = false
			}
			for field := range newSet {
				if !isContainedInNewOptionalAncestor(field, resourceDiff) {
					isContained = false
					break
				}
			}
			if isContained {
				continue
			}

			// Otherwise, all fields are breaking.
			for field := range newSet {
				messages = append(messages, fmt.Sprintf("Field `%s` within resource `%s` was added to exactly one of", field, resource))
			}
		}
	}
	return messages
}

func isComplexModification(newSet diff.FieldSet, resourceDiff diff.ResourceDiff) bool {
	for _, oldSet := range resourceDiff.FieldSets.Old.ExactlyOneOf {
		if len(newSet.Intersection(oldSet)) > 0 {
			return true
		}
	}
	return false
}

func isNewField(field string, diff diff.ResourceDiff) bool {
	fieldDiff, ok := diff.Fields[field]
	return !ok || fieldDiff.Old == nil
}

func isExistingFieldRequired(field string, diff diff.ResourceDiff) bool {
	fieldDiff, ok := diff.Fields[field]
	return ok && fieldDiff.Old != nil && fieldDiff.Old.Required
}

func isContainedInNewOptionalAncestor(field string, diff diff.ResourceDiff) bool {
	parts := strings.Split(field, ".")
	if len(parts) < 2 {
		return false
	}
	ancestorName := strings.Join(parts[:len(parts)-1], ".")
	ancestorDiff, ok := diff.Fields[ancestorName]
	if !ok {
		return false
	}

	isAncestorNew := ancestorDiff.Old == nil && ancestorDiff.New != nil
	isAncestorOptional := ancestorDiff.New != nil && ancestorDiff.New.Optional

	return isAncestorNew && isAncestorOptional
}
