package main

import (
	"strings"

	newProvider "google/provider/new/google-beta"
	oldProvider "google/provider/old/google-beta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Returns a map with resource names as keys and slices of changed field paths as values.
func changedResourceFields() map[string][]string {
	oldResourceMap := oldProvider.ResourceMap()
	newResourceMap := newProvider.ResourceMap()

	return resourceMapChanges(oldResourceMap, newResourceMap)
}

func resourceMapChanges(oldResourceMap, newResourceMap map[string]*schema.Resource) map[string][]string {
	changes := make(map[string][]string)
	for resourceName, newResource := range newResourceMap {
		if fields := changedFields(oldResourceMap[resourceName], newResource, nil); len(fields) > 0 {
			changes[resourceName] = fields
		}

	}
	return changes
}

func changedFields(oldResource, newResource *schema.Resource, path []string) []string {
	fields := make([]string, 0)
	for fieldName, newFieldSchema := range newResource.Schema {
		if fieldName == "project" {
			continue
		}
		if oldResource == nil {
			fields = append(fields, changedFieldsUnder(nil, newFieldSchema, append(path, fieldName))...)
		} else {
			fields = append(fields, changedFieldsUnder(oldResource.Schema[fieldName], newFieldSchema, append(path, fieldName))...)
		}
	}
	return fields
}

func changedFieldsUnder(oldFieldSchema, newFieldSchema *schema.Schema, path []string) []string {
	if newFieldSchema.Computed && !newFieldSchema.Optional {
		// Output only fields should not be included in missing test detection.
		return nil
	}
	if newFieldSchema.Elem != nil {
		if newFieldSchemaElem, ok := newFieldSchema.Elem.(*schema.Resource); ok {
			if oldFieldSchema != nil {
				return changedFields(oldFieldSchema.Elem.(*schema.Resource), newFieldSchemaElem, path)
			}
			return changedFields(nil, newFieldSchemaElem, path)
		}
	}
	if oldFieldSchema == nil || newFieldSchema.Type != oldFieldSchema.Type || newFieldSchema.Elem != oldFieldSchema.Elem {
		return []string{strings.Join(path, ".")}
	}
	return nil
}
