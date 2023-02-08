package main

import (
	newProvider "google/provider/new/google-beta"
	oldProvider "google/provider/old/google-beta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A map of field names with whether the field is covered as values.
type FieldCoverage map[string]any

// Returns a map with resource names as keys and field coverage maps as values.
func changedResourceFields() map[string]FieldCoverage {
	oldResourceMap := oldProvider.ResourceMap()
	newResourceMap := newProvider.ResourceMap()

	return resourceMapChanges(oldResourceMap, newResourceMap)
}

func resourceMapChanges(oldResourceMap, newResourceMap map[string]*schema.Resource) map[string]FieldCoverage {
	changes := make(map[string]FieldCoverage)
	for resourceName, newResource := range newResourceMap {
		if fields := changedFields(oldResourceMap[resourceName], newResource, false); len(fields) > 0 {
			changes[resourceName] = fields
		}

	}
	return changes
}

func changedFields(oldResource, newResource *schema.Resource, nested bool) FieldCoverage {
	fields := make(FieldCoverage)
	for fieldName, newFieldSchema := range newResource.Schema {
		if fieldName == "project" && !nested {
			// Skip checking the project field of resources since the provider automatically includes it.
			continue
		}
		var changed any
		if oldResource == nil {
			changed = changedFieldsUnder(nil, newFieldSchema)
		} else {
			changed = changedFieldsUnder(oldResource.Schema[fieldName], newFieldSchema)
		}
		if changed != nil {
			if coverage, ok := changed.(FieldCoverage); ok && len(coverage) == 0 {
				continue
			}
			fields[fieldName] = changed
		}
	}
	return fields
}

func changedFieldsUnder(oldFieldSchema, newFieldSchema *schema.Schema) any {
	if newFieldSchema.Computed && !newFieldSchema.Optional {
		// Output only fields should not be included in missing test detection.
		return nil
	}
	if newFieldSchema.Elem != nil {
		if newFieldSchemaElem, ok := newFieldSchema.Elem.(*schema.Resource); ok {
			if oldFieldSchema != nil {
				return changedFields(oldFieldSchema.Elem.(*schema.Resource), newFieldSchemaElem, true)
			}
			return changedFields(nil, newFieldSchemaElem, true)
		}
	}
	if oldFieldSchema == nil || newFieldSchema.Type != oldFieldSchema.Type || newFieldSchema.Elem != oldFieldSchema.Elem {
		return false
	}
	return nil
}
