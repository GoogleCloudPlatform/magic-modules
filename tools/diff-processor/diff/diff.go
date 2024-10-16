package diff

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SchemaDiff is a nested map with resource names as top-level keys.
// It includes diffs for all resources in the old and new providers.
type SchemaDiff map[string]ResourceDiff

type ResourceDiff struct {
	ResourceConfig ResourceConfigDiff
	Fields         map[string]FieldDiff
}

type ResourceConfigDiff struct {
	Old       *schema.Resource
	New       *schema.Resource
	Conflicts *FieldConflictSetsDiff // merged conflict set diffs for the entire resource
}

type FieldDiff struct {
	Old       *schema.Schema
	New       *schema.Schema
	Conflicts *FieldConflictSetsDiff // diffs in conflict sets (ConflictsWith, ExactlyOneOf, etc.) for *this* field
}

func ComputeSchemaDiff(oldResourceMap, newResourceMap map[string]*schema.Resource) SchemaDiff {
	schemaDiff := make(SchemaDiff)
	for resource := range union(oldResourceMap, newResourceMap) {
		// Compute diff between old and new resources and fields.
		// TODO: add support for computing diff between resource configs, not just whether the
		// resource was added/removed. b/300114839
		resourceDiff := ResourceDiff{}
		var flattenedOldSchema map[string]*schema.Schema
		if oldResource, ok := oldResourceMap[resource]; ok {
			flattenedOldSchema = flattenSchema("", oldResource.Schema)
			resourceDiff.ResourceConfig.Old = &schema.Resource{}
		}

		var flattenedNewSchema map[string]*schema.Schema
		if newResource, ok := newResourceMap[resource]; ok {
			flattenedNewSchema = flattenSchema("", newResource.Schema)
			resourceDiff.ResourceConfig.New = &schema.Resource{}
		}

		resourceDiff.Fields = make(map[string]FieldDiff)
		for key := range union(flattenedOldSchema, flattenedNewSchema) {
			oldField := flattenedOldSchema[key]
			newField := flattenedNewSchema[key]
			fd := diffFields(oldField, newField)
			if fd != nil {
				resourceDiff.Fields[key] = *fd
				if fd.Conflicts != nil {
					if resourceDiff.ResourceConfig.Conflicts == nil {
						resourceDiff.ResourceConfig.Conflicts = fd.Conflicts
					} else {
						resourceDiff.ResourceConfig.Conflicts.Merge(fd.Conflicts)
					}
				}
			}
		}
		if len(resourceDiff.Fields) > 0 || !cmp.Equal(resourceDiff.ResourceConfig.Old, resourceDiff.ResourceConfig.New) {
			schemaDiff[resource] = resourceDiff
		}
	}
	return schemaDiff
}

func flattenSchema(parentKey string, schemaObj map[string]*schema.Schema) map[string]*schema.Schema {
	flattened := make(map[string]*schema.Schema)

	if parentKey != "" {
		parentKey += "."
	}

	for fieldName, field := range schemaObj {
		key := parentKey + fieldName
		flattened[key] = field
		childResource, hasNestedFields := field.Elem.(*schema.Resource)
		if field.Elem != nil && hasNestedFields {
			for childKey, childField := range flattenSchema(key, childResource.Schema) {
				flattened[childKey] = childField
			}
		}
	}

	return flattened
}

func diffFields(oldField, newField *schema.Schema) *FieldDiff {
	// If either field is nil, it is changed; if both are nil (which should never happen) it's not
	if oldField == nil && newField == nil {
		return nil
	}

	diff := &FieldDiff{
		Old: oldField,
		New: newField,
	}

	if oldField == nil {
		diff.Conflicts = diffFieldConflictSets(
			nil,
			makeFieldConflictSetsFromSchema(newField),
		)
		return diff
	}

	if newField == nil {
		diff.Conflicts = diffFieldConflictSets(
			makeFieldConflictSetsFromSchema(oldField),
			nil,
		)
		return diff
	}

	hasDiff := false

	if basicSchemaChanged(oldField, newField) {
		hasDiff = true
	}

	diff.Conflicts = diffFieldConflictSets(
		makeFieldConflictSetsFromSchema(oldField),
		makeFieldConflictSetsFromSchema(newField),
	)
	if diff.Conflicts != nil {
		hasDiff = true
	}

	if elemChanged(oldField, newField) {
		hasDiff = true
	}

	if funcsChanged(oldField, newField) {
		hasDiff = true
	}

	if !hasDiff {
		return nil
	}
	return diff
}

// Check if any basic Schema struct fields have changed.
// https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.24.0/helper/schema/schema.go#L44
func basicSchemaChanged(oldField, newField *schema.Schema) bool {
	if oldField.Type != newField.Type {
		return true
	}
	if oldField.ConfigMode != newField.ConfigMode {
		return true
	}
	if oldField.Required != newField.Required {
		return true
	}
	if oldField.Optional != newField.Optional {
		return true
	}
	if oldField.Computed != newField.Computed {
		return true
	}
	if oldField.ForceNew != newField.ForceNew {
		return true
	}
	if oldField.DiffSuppressOnRefresh != newField.DiffSuppressOnRefresh {
		return true
	}
	if oldField.Default != newField.Default {
		return true
	}
	if oldField.Description != newField.Description {
		return true
	}
	if oldField.InputDefault != newField.InputDefault {
		return true
	}
	if oldField.MaxItems != newField.MaxItems {
		return true
	}
	if oldField.MinItems != newField.MinItems {
		return true
	}
	if oldField.Deprecated != newField.Deprecated {
		return true
	}
	if oldField.Sensitive != newField.Sensitive {
		return true
	}
	return false
}

// Check if Elem changed (unless old and new both represent nested fields)
func elemChanged(oldField, newField *schema.Schema) bool {
	if (oldField.Elem == nil && newField.Elem != nil) || (oldField.Elem != nil && newField.Elem == nil) {
		return true
	}
	if oldField.Elem != nil && newField.Elem != nil {
		// At this point new/old Elems are either schema.Schema or schema.Resource.
		// If both are schema.Resource we don't need to do anything. Diffs on subfields
		// are handled separately.
		_, oldIsResource := oldField.Elem.(*schema.Resource)
		_, newIsResource := newField.Elem.(*schema.Resource)

		if (oldIsResource && !newIsResource) || (!oldIsResource && newIsResource) {
			return true
		}
		if !oldIsResource && !newIsResource {
			if fd := diffFields(oldField.Elem.(*schema.Schema), newField.Elem.(*schema.Schema)); fd != nil {
				return true
			}
		}
	}
	return false
}

// Check if any Schema struct fields that are functions have changed
func funcsChanged(oldField, newField *schema.Schema) bool {
	if funcChanged(oldField.DiffSuppressFunc, newField.DiffSuppressFunc) {
		return true
	}
	if funcChanged(oldField.DefaultFunc, newField.DefaultFunc) {
		return true
	}
	if funcChanged(oldField.StateFunc, newField.StateFunc) {
		return true
	}
	if funcChanged(oldField.Set, newField.Set) {
		return true
	}
	if funcChanged(oldField.ValidateFunc, newField.ValidateFunc) {
		return true
	}
	if funcChanged(oldField.ValidateDiagFunc, newField.ValidateDiagFunc) {
		return true
	}
	return false
}

func funcChanged(oldFunc, newFunc interface{}) bool {
	// If it changed to/from nil, it changed
	oldFuncIsNil := reflect.ValueOf(oldFunc).IsNil()
	newFuncIsNil := reflect.ValueOf(newFunc).IsNil()
	if (oldFuncIsNil && !newFuncIsNil) || (!oldFuncIsNil && newFuncIsNil) {
		return true
	}

	// If a func is set before and after we don't currently have a way to reliably
	// determine whether the function changed, so we assume that it has not changed.
	// b/300157205
	return false
}
