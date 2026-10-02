package fwutils

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EphemeralRequestField describes the API names and input constraints independently
// of response flattening. Defaults are JSON so false, zero, and empty strings are
// distinguishable from an absent default.
type EphemeralRequestField struct {
	ApiName     string
	AttrType    attr.Type
	Required    bool
	Output      bool
	DefaultJSON string
	MinSize     int
	MaxSize     int
	Enum        []string
	Pattern     string
	Fields      map[string]EphemeralRequestField
	Element     *EphemeralRequestField
}

// ExpandEphemeralValue returns JSON-compatible request data and the input value
// to retain in the result. Nulls are omitted by the caller, while explicit empty
// values are preserved. Unknown values are diagnosed before making an API call.
func ExpandEphemeralValue(ctx context.Context, value attr.Value, field EphemeralRequestField) (interface{}, attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	fail := func(message string) (interface{}, attr.Value, diag.Diagnostics) {
		diags.AddError("Invalid ephemeral input", fmt.Sprintf("Field %q: %s", field.ApiName, message))
		return nil, value, diags
	}
	if value.IsUnknown() {
		return fail("value must be known before opening the resource")
	}
	if value.IsNull() && field.DefaultJSON != "" {
		var raw interface{}
		decoder := json.NewDecoder(strings.NewReader(field.DefaultJSON))
		decoder.UseNumber()
		if err := decoder.Decode(&raw); err != nil {
			return fail("invalid generated default")
		}
		if number, ok := raw.(json.Number); ok {
			raw = number.String()
		}
		value, diags = FlattenEphemeralValue(raw, EphemeralFieldMapping{AttrType: field.AttrType})
		if diags.HasError() {
			return nil, value, diags
		}
	}
	if value.IsNull() {
		if field.Required {
			return fail("a value is required")
		}
		return nil, value, diags
	}
	checkSize := func(size int) bool {
		return size >= field.MinSize && (field.MaxSize == 0 || size <= field.MaxSize)
	}
	switch v := value.(type) {
	case types.String:
		s := v.ValueString()
		if len(field.Enum) > 0 && !slices.Contains(field.Enum, s) {
			return fail("value is not in the allowed enum")
		}
		if field.Pattern != "" {
			r, err := regexp.Compile(field.Pattern)
			if err != nil || !r.MatchString(s) {
				return fail("value does not match the required pattern")
			}
		}
		return s, value, diags
	case types.Bool:
		return v.ValueBool(), value, diags
	case types.Int64:
		return v.ValueInt64(), value, diags
	case types.Float64:
		return v.ValueFloat64(), value, diags
	case types.Object:
		body := make(map[string]interface{})
		attrs := v.Attributes()
		for name, child := range field.Fields {
			if child.Output {
				continue
			}
			input, ok := attrs[name]
			if !ok {
				return fail("object attribute missing from generated schema")
			}
			raw, normalized, d := ExpandEphemeralValue(ctx, input, child)
			diags.Append(d...)
			if d.HasError() {
				continue
			}
			attrs[name] = normalized
			if !normalized.IsNull() {
				body[child.ApiName] = raw
			}
		}
		if diags.HasError() {
			return nil, value, diags
		}
		normalized, d := types.ObjectValue(v.AttributeTypes(ctx), attrs)
		diags.Append(d...)
		return body, normalized, diags
	case types.List, types.Set:
		var elements []attr.Value
		var elementType attr.Type
		switch collection := value.(type) {
		case types.List:
			elements, elementType = collection.Elements(), collection.ElementType(ctx)
		case types.Set:
			elements, elementType = collection.Elements(), collection.ElementType(ctx)
		}
		if !checkSize(len(elements)) {
			return fail("collection size is outside the allowed bounds")
		}
		child := EphemeralRequestField{AttrType: elementType}
		if field.Element != nil {
			child = *field.Element
		}
		body := make([]interface{}, 0, len(elements))
		for i, input := range elements {
			if input.IsNull() {
				return fail("collection elements must not be null")
			}
			raw, normalized, d := ExpandEphemeralValue(ctx, input, child)
			diags.Append(d...)
			elements[i] = normalized
			body = append(body, raw)
		}
		if diags.HasError() {
			return nil, value, diags
		}
		if _, ok := value.(types.List); ok {
			normalized, d := types.ListValue(elementType, elements)
			diags.Append(d...)
			return body, normalized, diags
		}
		normalized, d := types.SetValue(elementType, elements)
		diags.Append(d...)
		return body, normalized, diags
	case types.Map:
		elements := v.Elements()
		if !checkSize(len(elements)) {
			return fail("map size is outside the allowed bounds")
		}
		child := EphemeralRequestField{AttrType: v.ElementType(ctx)}
		if field.Element != nil {
			child = *field.Element
		}
		body := make(map[string]interface{}, len(elements))
		for key, input := range elements {
			raw, normalized, d := ExpandEphemeralValue(ctx, input, child)
			diags.Append(d...)
			if d.HasError() {
				continue
			}
			elements[key] = normalized
			if !normalized.IsNull() {
				body[key] = raw
			}
		}
		if diags.HasError() {
			return nil, value, diags
		}
		normalized, d := types.MapValue(v.ElementType(ctx), elements)
		diags.Append(d...)
		return body, normalized, diags
	default:
		return fail("unsupported Framework value type")
	}
}
