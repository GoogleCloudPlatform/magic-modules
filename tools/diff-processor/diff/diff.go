package diff

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/maps"
)

// SchemaDiff is a nested map with resource names as top-level keys.
type SchemaDiff map[string]ResourceDiff

type ResourceDiff struct {
	ResourceConfig ResourceConfigDiff
	Fields         map[string]FieldDiff
}

type ResourceConfigDiff struct {
	Old *schema.Resource
	New *schema.Resource
}

type FieldDiff struct {
	Changed map[string]bool
	Old     *schema.Schema
	New     *schema.Schema
}

func ComputeSchemaDiff(oldResourceMap, newResourceMap map[string]*schema.Resource) SchemaDiff {
	schemaDiff := make(SchemaDiff)
	for resource, _ := range union(maps.Keys(oldResourceMap), maps.Keys(newResourceMap)) {
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
		for key, _ := range union(maps.Keys(flattenedOldSchema), maps.Keys(flattenedNewSchema)) {
			oldField := flattenedOldSchema[key]
			newField := flattenedNewSchema[key]
			if hasChanges, changed := fieldChanged(oldField, newField); hasChanges {
				resourceDiff.Fields[key] = FieldDiff{
					Changed: changed,
					Old:     oldField,
					New:     newField,
				}
			}
		}
		if len(resourceDiff.Fields) > 0 || !cmp.Equal(resourceDiff.ResourceConfig.Old, resourceDiff.ResourceConfig.New) {
			schemaDiff[resource] = resourceDiff
		}
	}
	return schemaDiff
}

func union(keys1, keys2 []string) map[string]struct{} {
	allKeys := make(map[string]struct{})
	for _, key := range keys1 {
		allKeys[key] = struct{}{}
	}
	for _, key := range keys2 {
		allKeys[key] = struct{}{}
	}
	return allKeys
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

func fieldChanged(oldField, newField *schema.Schema) (bool, map[string]bool) {
	// If either field is nil, it is changed; if both are nil (which should never happen) it's not
	if oldField == nil && newField == nil {
		return false, nil
	}
	if oldField == nil || newField == nil {
		return true, nil
	}
	changed := map[string]bool{}
	// Check if any basic Schema struct fields have changed.
	// https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.24.0/helper/schema/schema.go#L44
	if oldField.Type != newField.Type {
		changed["Type"] = true
	}
	if oldField.ConfigMode != newField.ConfigMode {
		changed["ConfigMode"] = true
	}
	if oldField.Required != newField.Required {
		changed["Required"] = true
	}
	if oldField.Optional != newField.Optional {
		changed["Optional"] = true
	}
	if oldField.Computed != newField.Computed {
		changed["Computed"] = true
	}
	if oldField.ForceNew != newField.ForceNew {
		changed["ForceNew"] = true
	}
	if oldField.DiffSuppressOnRefresh != newField.DiffSuppressOnRefresh {
		changed["DiffSuppressOnRefresh"] = true
	}
	if oldField.Default != newField.Default {
		changed["Default"] = true
	}
	if oldField.Description != newField.Description {
		changed["Description"] = true
	}
	if oldField.InputDefault != newField.InputDefault {
		changed["InputDefault"] = true
	}
	if oldField.MaxItems != newField.MaxItems {
		changed["MaxItems"] = true
	}
	if oldField.MinItems != newField.MinItems {
		changed["MinItems"] = true
	}
	if oldField.Deprecated != newField.Deprecated {
		changed["Deprecated"] = true
	}
	if oldField.Sensitive != newField.Sensitive {
		changed["Sensitive"] = true
	}

	// Compare slices
	less := func(a, b string) bool { return a < b }

	if (len(oldField.ConflictsWith) > 0 || len(newField.ConflictsWith) > 0) && !cmp.Equal(oldField.ConflictsWith, newField.ConflictsWith, cmpopts.SortSlices(less)) {
		changed["ConflictsWith"] = true
	}

	if (len(oldField.ExactlyOneOf) > 0 || len(newField.ExactlyOneOf) > 0) && !cmp.Equal(oldField.ExactlyOneOf, newField.ExactlyOneOf, cmpopts.SortSlices(less)) {
		changed["ExactlyOneOf"] = true
	}

	if (len(oldField.AtLeastOneOf) > 0 || len(newField.AtLeastOneOf) > 0) && !cmp.Equal(oldField.AtLeastOneOf, newField.AtLeastOneOf, cmpopts.SortSlices(less)) {
		changed["AtLeastOneOf"] = true
	}

	if (len(oldField.RequiredWith) > 0 || len(newField.RequiredWith) > 0) && !cmp.Equal(oldField.RequiredWith, newField.RequiredWith, cmpopts.SortSlices(less)) {
		changed["RequiredWith"] = true
	}

	// Check if Elem changed (unless old and new both represent nested fields)
	if (oldField.Elem == nil && newField.Elem != nil) || (oldField.Elem != nil && newField.Elem == nil) {
		changed["Elem"] = true
	}
	if oldField.Elem != nil && newField.Elem != nil {
		// At this point new/old Elems are either schema.Schema or schema.Resource.
		// If both are schema.Resource we don't need to do anything. Diffs on subfields
		// are handled separately.
		_, oldIsResource := oldField.Elem.(*schema.Resource)
		_, newIsResource := newField.Elem.(*schema.Resource)

		if (oldIsResource && !newIsResource) || (!oldIsResource && newIsResource) {
			changed["Elem"] = true
		}
		if !oldIsResource && !newIsResource {
			if hasChanges, elemChanged := fieldChanged(oldField.Elem.(*schema.Schema), newField.Elem.(*schema.Schema)); hasChanges {
				for k := range elemChanged {
					changed[fmt.Sprintf("Elem.%s", k)] = true
				}
			}
		}
	}

	// Check if any Schema struct fields that are functions have changed
	if funcChanged(oldField.DiffSuppressFunc, newField.DiffSuppressFunc) {
		changed["DiffSuppressFunc"] = true
	}
	if funcChanged(oldField.DefaultFunc, newField.DefaultFunc) {
		changed["DefaultFunc"] = true
	}
	if funcChanged(oldField.StateFunc, newField.StateFunc) {
		changed["StateFunc"] = true
	}
	if funcChanged(oldField.Set, newField.Set) {
		changed["Set"] = true
	}
	if funcChanged(oldField.ValidateFunc, newField.ValidateFunc) {
		changed["ValidateFunc"] = true
	}
	if funcChanged(oldField.ValidateDiagFunc, newField.ValidateDiagFunc) {
		changed["ValidateDiagFunc"] = true
	}

	if len(changed) > 0 {
		return true, changed
	}
	return false, nil
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
