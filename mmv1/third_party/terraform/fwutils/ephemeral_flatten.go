package fwutils

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EphemeralFieldMapping describes how to convert one field of a raw JSON API
// response (as returned by transport_tpg.SendRequest: map[string]interface{},
// []interface{}, string, float64, bool, or nil) into a Terraform attr.Value for
// a generated ephemeral resource, including nested objects/lists/sets/maps at
// arbitrary depth. Generated per-property by the ephemeral resource template;
// the recursive walk itself lives once in FlattenEphemeralValue below.
type EphemeralFieldMapping struct {
	// ApiName is this field's key in the raw API response map. Unused for the
	// synthetic mapping passed as the element/value mapping of a List/Set/Map.
	ApiName string
	// AttrType is the Terraform attribute type for this field: a scalar type,
	// or types.ObjectType/ListType/SetType/MapType for structured data.
	AttrType attr.Type
	// Fields describes an object-shaped value's sub-fields (when AttrType is
	// types.ObjectType), or the element/value mapping to apply to each item
	// (when AttrType is types.ListType/SetType/MapType). Keyed by Terraform
	// attribute name; ignored for scalar AttrTypes.
	Fields map[string]EphemeralFieldMapping
	// Transform runs before type conversion, including for a missing API value.
	// It returns API-shaped data, not SDK ResourceData or Framework values.
	Transform func(interface{}) (interface{}, error)
}

// FlattenEphemeralValue converts a raw JSON value into an attr.Value matching
// mapping.AttrType, recursing through mapping.Fields for nested object/list/
// set/map shapes. The returned attr.Value's concrete type always matches
// mapping.AttrType's value type (e.g. types.StringType -> types.String), so
// callers can safely type-assert it into a generated model struct field.
func FlattenEphemeralValue(raw interface{}, mapping EphemeralFieldMapping) (attr.Value, diag.Diagnostics) {
	if mapping.Transform != nil {
		var err error
		raw, err = mapping.Transform(raw)
		if err != nil {
			var diags diag.Diagnostics
			// Transformation errors may include secret material. Identify the field,
			// but do not put the raw value or arbitrary error text in diagnostics.
			diags.AddError("Error transforming API response", fmt.Sprintf("Could not transform API field %q.", mapping.ApiName))
			return nil, diags
		}
	}
	switch t := mapping.AttrType.(type) {
	case types.ObjectType:
		return flattenEphemeralObject(raw, t, mapping.Fields)
	case types.ListType:
		return flattenEphemeralList(raw, t, mapping.Fields)
	case types.SetType:
		return flattenEphemeralSet(raw, t, mapping.Fields)
	case types.MapType:
		return flattenEphemeralMap(raw, t, mapping.Fields)
	default:
		return flattenEphemeralScalar(raw, mapping.AttrType)
	}
}

func flattenEphemeralObject(raw interface{}, t types.ObjectType, fields map[string]EphemeralFieldMapping) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	m, ok := raw.(map[string]interface{})
	if !ok {
		if raw != nil {
			return unexpectedEphemeralType(raw, t)
		}
		return types.ObjectNull(t.AttrTypes), diags
	}

	attrs := make(map[string]attr.Value, len(fields))
	for attrName, fieldMapping := range fields {
		val, d := FlattenEphemeralValue(m[fieldMapping.ApiName], fieldMapping)
		diags.Append(d...)
		attrs[attrName] = val
	}
	if diags.HasError() {
		return nil, diags
	}

	objVal, d := types.ObjectValue(t.AttrTypes, attrs)
	diags.Append(d...)
	return objVal, diags
}

func flattenEphemeralList(raw interface{}, t types.ListType, elemMapping map[string]EphemeralFieldMapping) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	arr, ok := raw.([]interface{})
	if !ok {
		if raw != nil {
			return unexpectedEphemeralType(raw, t)
		}
		return types.ListNull(t.ElemType), diags
	}

	mapping := EphemeralFieldMapping{AttrType: t.ElemType, Fields: elemMapping}
	values := make([]attr.Value, 0, len(arr))
	for _, item := range arr {
		val, d := FlattenEphemeralValue(item, mapping)
		diags.Append(d...)
		values = append(values, val)
	}
	if diags.HasError() {
		return nil, diags
	}

	listVal, d := types.ListValue(t.ElemType, values)
	diags.Append(d...)
	return listVal, diags
}

func flattenEphemeralSet(raw interface{}, t types.SetType, elemMapping map[string]EphemeralFieldMapping) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	arr, ok := raw.([]interface{})
	if !ok {
		if raw != nil {
			return unexpectedEphemeralType(raw, t)
		}
		return types.SetNull(t.ElemType), diags
	}

	mapping := EphemeralFieldMapping{AttrType: t.ElemType, Fields: elemMapping}
	values := make([]attr.Value, 0, len(arr))
	for _, item := range arr {
		val, d := FlattenEphemeralValue(item, mapping)
		diags.Append(d...)
		values = append(values, val)
	}
	if diags.HasError() {
		return nil, diags
	}

	setVal, d := types.SetValue(t.ElemType, values)
	diags.Append(d...)
	return setVal, diags
}

func flattenEphemeralMap(raw interface{}, t types.MapType, valueMapping map[string]EphemeralFieldMapping) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	m, ok := raw.(map[string]interface{})
	if !ok {
		if raw != nil {
			return unexpectedEphemeralType(raw, t)
		}
		return types.MapNull(t.ElemType), diags
	}

	mapping := EphemeralFieldMapping{AttrType: t.ElemType, Fields: valueMapping}
	values := make(map[string]attr.Value, len(m))
	for k, v := range m {
		val, d := FlattenEphemeralValue(v, mapping)
		diags.Append(d...)
		values[k] = val
	}
	if diags.HasError() {
		return nil, diags
	}

	mapVal, d := types.MapValue(t.ElemType, values)
	diags.Append(d...)
	return mapVal, diags
}

func flattenEphemeralScalar(raw interface{}, t attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch t {
	case types.StringType:
		switch v := raw.(type) {
		case nil:
			return types.StringNull(), diags
		case string:
			return types.StringValue(v), diags
		}
	case types.BoolType:
		switch v := raw.(type) {
		case nil:
			return types.BoolNull(), diags
		case bool:
			return types.BoolValue(v), diags
		}
	case types.Int64Type:
		switch v := raw.(type) {
		case nil:
			return types.Int64Null(), diags
		case float64:
			return types.Int64Value(int64(v)), diags
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(i), diags
			}
		}
	case types.Float64Type:
		switch v := raw.(type) {
		case nil:
			return types.Float64Null(), diags
		case float64:
			return types.Float64Value(v), diags
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return types.Float64Value(f), diags
			}
		}
	}

	return unexpectedEphemeralType(raw, t)
}

func unexpectedEphemeralType(raw interface{}, t attr.Type) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	diags.AddError(
		"Unexpected API response type",
		fmt.Sprintf("Could not convert API value of type %T into attribute type %s.", raw, t.String()),
	)
	return nil, diags
}
