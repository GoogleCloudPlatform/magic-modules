package main

import (
	newProvider "google/provider/new/google-beta/provider"
	oldProvider "google/provider/old/google-beta/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceChanges is a nested map with field names as keys and Field objects
// as bottom-level values.
// Fields are assumed not to be covered until detected in a test.
type ResourceChanges map[string]any

type Field struct {
	// Added is true when the field is newly added between oldProvider and newProvider.
	Added bool
	// Changed is true when the field type has changed between oldProvider and newProvider.
	Changed bool
	// Tested is true when a test has been found that includes the field.
	Tested bool
}

// Returns a map with resource names as keys and field coverage maps as values.
func changedResourceFields() map[string]ResourceChanges {
	oldResourceMap := oldProvider.ResourceMap()
	newResourceMap := newProvider.ResourceMap()

	return resourceMapChanges(oldResourceMap, newResourceMap)
}

func resourceMapChanges(oldResourceMap, newResourceMap map[string]*schema.Resource) map[string]ResourceChanges {
	changes := make(map[string]ResourceChanges)
	for resourceName, newResource := range newResourceMap {
		if resourceName == "google_compute_instance_from_template" || resourceName == "google_compute_instance_from_machine_image" {
			// This resource is skipped because its changes can be covered by google_compute_instance.
			continue
		}
		if fields := changedFields(oldResourceMap[resourceName], newResource, false); len(fields) > 0 {
			changes[resourceName] = fields
		}

	}
	return changes
}

func changedFields(oldResource, newResource *schema.Resource, nested bool) ResourceChanges {
	fields := make(ResourceChanges)
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
			if resourceChanges, ok := changed.(ResourceChanges); ok && len(resourceChanges) == 0 {
				continue
			}
			fields[fieldName] = changed
		}
	}
	return fields
}

// Return a Field for changed non-nested fields or a ResourceChanges map for changed nested fields.
// Return nil if the schemas are identical.
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
	if oldFieldSchema == nil {
		return &Field{Added: true}
	}
	if oldFieldSchema.Type != newFieldSchema.Type {
		return &Field{Changed: true}
	}
	// TODO(trodge) handle the case where something under Elem changed but Elem is not a schema.Resource
	return nil
}
