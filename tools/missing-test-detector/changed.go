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
		if oldResource, ok := oldResourceMap[resourceName]; !ok {
			changes[resourceName] = allFields(newResource, nil)
		} else if fields := changedFields(oldResource, newResource, nil); len(fields) > 0 {
			changes[resourceName] = fields
		}

	}
	return changes
}

func allFields(resource *schema.Resource, path []string) []string {
	fields := make([]string, 0)
	for fieldName, fieldSchema := range resource.Schema {
		if fieldName == "project" {
			continue
		}
		fields = append(fields, allFieldsUnder(fieldSchema, append(path, fieldName))...)
	}
	return fields
}

func allFieldsUnder(fieldSchema *schema.Schema, path []string) []string {
	if fieldSchema.Computed && !fieldSchema.Optional {
		// Output only fields should not be included in missing test detection.
		return nil
	}
	if fieldSchema.Elem != nil {
		if fieldSchemaElem, ok := fieldSchema.Elem.(*schema.Resource); ok {
			// Field is a nested object.
			return allFields(fieldSchemaElem, path)
		}
	}
	// Field is a scalar, map, or list.
	return []string{strings.Join(path, ".")}
}

func changedFields(oldResource, newResource *schema.Resource, path []string) []string {
	fields := make([]string, 0)
	for fieldName, newFieldSchema := range newResource.Schema {
		if fieldName == "project" {
			continue
		}
		if oldFieldSchema, ok := oldResource.Schema[fieldName]; ok {
			fields = append(fields, changedFieldsUnder(oldFieldSchema, newFieldSchema, append(path, fieldName))...)
		} else {
			fields = append(fields, allFieldsUnder(newFieldSchema, append(path, fieldName))...)
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
		if oldFieldSchema.Elem == nil {
			return allFieldsUnder(newFieldSchema, path)
		}
		if newFieldSchemaElem, ok := newFieldSchema.Elem.(*schema.Resource); ok {
			if oldFieldSchemaElem, ok := oldFieldSchema.Elem.(*schema.Resource); ok {
				return changedFields(oldFieldSchemaElem, newFieldSchemaElem, path)
			}
			return allFields(newFieldSchemaElem, path)
		}
	}
	if newFieldSchema.Type != oldFieldSchema.Type {
		return []string{strings.Join(path, ".")}
	}
	return nil
}
