package diff

import (
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SchemaDiff is a nested map with resource names as top-level keys.
type SchemaDiff map[string]ResourceDiff

type ResourceDiffInterface interface {
	IsNewResource() bool
	IsFieldInNewNestedStructure(fieldPath string) bool
}

type ResourceDiff struct {
	ResourceConfig  ResourceConfigDiff
	FlattenedSchema FlattenedSchemaRaw
	Fields          map[string]FieldDiff
	FieldSets       ResourceFieldSetsDiff
}

type ResourceFieldSetsDiff struct {
	Old ResourceFieldSets
	New ResourceFieldSets
}

type ResourceFieldSets struct {
	ConflictsWith map[string]FieldSet
	ExactlyOneOf  map[string]FieldSet
	AtLeastOneOf  map[string]FieldSet
	RequiredWith  map[string]FieldSet
}

type FieldSet map[string]struct{}

type ResourceConfigDiff struct {
	Old *schema.Resource
	New *schema.Resource
}

type FieldDiff struct {
	Old *schema.Schema
	New *schema.Schema
}

type FlattenedSchemaRaw struct {
	Old map[string]*schema.Schema
	New map[string]*schema.Schema
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
			resourceDiff.FlattenedSchema.Old = flattenedOldSchema
			resourceDiff.ResourceConfig.Old = &schema.Resource{}
		}

		var flattenedNewSchema map[string]*schema.Schema
		if newResource, ok := newResourceMap[resource]; ok {
			flattenedNewSchema = flattenSchema("", newResource.Schema)
			resourceDiff.FlattenedSchema.New = flattenedNewSchema
			resourceDiff.ResourceConfig.New = &schema.Resource{}
		}

		resourceDiff.Fields = make(map[string]FieldDiff)
		for key := range union(flattenedOldSchema, flattenedNewSchema) {
			oldField := flattenedOldSchema[key]
			newField := flattenedNewSchema[key]
			if fieldDiff, fieldSetsDiff, changed := diffFields(oldField, newField, key); changed {
				resourceDiff.Fields[key] = fieldDiff
				resourceDiff.FieldSets = mergeFieldSetsDiff(resourceDiff.FieldSets, fieldSetsDiff)
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

func diffFields(oldField, newField *schema.Schema, fieldName string) (FieldDiff, ResourceFieldSetsDiff, bool) {
	// If either field is nil, it is changed; if both are nil (which should never happen) it's not
	if oldField == nil && newField == nil {
		return FieldDiff{}, ResourceFieldSetsDiff{}, false
	}

	oldFieldSets := fieldSets(oldField, fieldName)
	newFieldSets := fieldSets(newField, fieldName)

	fieldDiff := FieldDiff{
		Old: oldField,
		New: newField,
	}
	fieldSetsDiff := ResourceFieldSetsDiff{
		Old: oldFieldSets,
		New: newFieldSets,
	}
	if oldField == nil || newField == nil {
		return fieldDiff, fieldSetsDiff, true
	}
	// Check if any basic Schema struct fields have changed.
	// https://github.com/hashicorp/terraform-plugin-sdk/blob/v2.24.0/helper/schema/schema.go#L44
	if basicSchemaChanged(oldField, newField) {
		return fieldDiff, fieldSetsDiff, true
	}

	if !cmp.Equal(oldFieldSets, newFieldSets) {
		return fieldDiff, fieldSetsDiff, true
	}

	if elemChanged(oldField, newField) {
		return fieldDiff, fieldSetsDiff, true
	}

	if funcsChanged(oldField, newField) {
		return fieldDiff, fieldSetsDiff, true
	}

	return FieldDiff{}, ResourceFieldSetsDiff{}, false
}

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

func fieldSets(field *schema.Schema, fieldName string) ResourceFieldSets {
	if field == nil {
		return ResourceFieldSets{}
	}
	var conflictsWith, exactlyOneOf, atLeastOneOf, requiredWith map[string]FieldSet
	if len(field.ConflictsWith) > 0 {
		set := sliceToSetRemoveZeroPadding(append(field.ConflictsWith, fieldName))
		conflictsWith = map[string]FieldSet{
			setKey(set): set,
		}
	}
	if len(field.ExactlyOneOf) > 0 {
		set := sliceToSetRemoveZeroPadding(append(field.ExactlyOneOf, fieldName))
		exactlyOneOf = map[string]FieldSet{
			setKey(set): set,
		}
	}
	if len(field.AtLeastOneOf) > 0 {
		set := sliceToSetRemoveZeroPadding(append(field.AtLeastOneOf, fieldName))
		atLeastOneOf = map[string]FieldSet{
			setKey(set): set,
		}
	}
	if len(field.RequiredWith) > 0 {
		set := sliceToSetRemoveZeroPadding(append(field.RequiredWith, fieldName))
		requiredWith = map[string]FieldSet{
			setKey(set): set,
		}
	}
	return ResourceFieldSets{
		ConflictsWith: conflictsWith,
		ExactlyOneOf:  exactlyOneOf,
		AtLeastOneOf:  atLeastOneOf,
		RequiredWith:  requiredWith,
	}
}

func elemChanged(oldField, newField *schema.Schema) bool {
	// Check if Elem changed (unless old and new both represent nested fields)
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
			if _, _, changed := diffFields(oldField.Elem.(*schema.Schema), newField.Elem.(*schema.Schema), ""); changed {
				return true
			}
		}
	}
	return false
}

func funcsChanged(oldField, newField *schema.Schema) bool {
	// Check if any Schema struct fields that are functions have changed
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

func mergeFieldSetsDiff(allFields ResourceFieldSetsDiff, currentField ResourceFieldSetsDiff) ResourceFieldSetsDiff {
	allFields.Old = mergeResourceFieldSets(allFields.Old, currentField.Old)
	allFields.New = mergeResourceFieldSets(allFields.New, currentField.New)
	return allFields
}

func mergeResourceFieldSets(allFields ResourceFieldSets, currentField ResourceFieldSets) ResourceFieldSets {
	allFields.ConflictsWith = mergeFieldSets(allFields.ConflictsWith, currentField.ConflictsWith)
	allFields.ExactlyOneOf = mergeFieldSets(allFields.ExactlyOneOf, currentField.ExactlyOneOf)
	allFields.AtLeastOneOf = mergeFieldSets(allFields.AtLeastOneOf, currentField.AtLeastOneOf)
	allFields.RequiredWith = mergeFieldSets(allFields.RequiredWith, currentField.RequiredWith)
	return allFields
}

func mergeFieldSets(allFields, currentField map[string]FieldSet) map[string]FieldSet {
	if allFields == nil {
		allFields = make(map[string]FieldSet)
	}
	for key, fieldSet := range currentField {
		allFields[key] = fieldSet
	}
	if len(allFields) == 0 {
		return nil
	}
	return allFields
}

func setKey(set FieldSet) string {
	slice := setToSortedSlice(set)
	return strings.Join(slice, ",")
}

func (rd ResourceDiff) IsNewResource() bool {
	rcd := rd.ResourceConfig
	if rcd.Old == nil && rcd.New != nil {
		return true
	}
	return false
}

// IsFieldInNewNestedStructure determines if a field is part of a completely new nested structure
func (rd ResourceDiff) IsFieldInNewNestedStructure(fieldPath string) bool {
	if rd.IsNewResource() {
		return true
	}

	// Get the parent path (everything before the last dot)
	lastDotIndex := strings.LastIndex(fieldPath, ".")
	if lastDotIndex == -1 {
		// No parent path (top-level field)
		return false
	}

	parentPath := fieldPath[:lastDotIndex]

	// Check if parent exists in new schema but not in old schema
	_, parentExistsInOld := rd.FlattenedSchema.Old[parentPath]
	_, parentExistsInNew := rd.FlattenedSchema.New[parentPath]

	return !parentExistsInOld && parentExistsInNew
}
