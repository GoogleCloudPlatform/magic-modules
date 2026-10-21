package fwutils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenEphemeralValue_Scalar(t *testing.T) {
	cases := []struct {
		name     string
		raw      interface{}
		attrType attr.Type
		want     attr.Value
	}{
		{"string", "hello", types.StringType, types.StringValue("hello")},
		{"string nil", nil, types.StringType, types.StringNull()},
		{"bool", true, types.BoolType, types.BoolValue(true)},
		{"int64 from float64", float64(42), types.Int64Type, types.Int64Value(42)},
		{"int64 from string", "42", types.Int64Type, types.Int64Value(42)},
		{"float64", float64(3.5), types.Float64Type, types.Float64Value(3.5)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, diags := FlattenEphemeralValue(tc.raw, EphemeralFieldMapping{AttrType: tc.attrType})
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if !got.Equal(tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestFlattenEphemeralValue_SingleNestedObject(t *testing.T) {
	// Mirrors the Phase 0 spike shape: effective_intelligence_config { intelligence_config, effective_edition }
	attrTypes := map[string]attr.Type{
		"intelligence_config": types.StringType,
		"effective_edition":   types.StringType,
	}
	mapping := EphemeralFieldMapping{
		AttrType: types.ObjectType{AttrTypes: attrTypes},
		Fields: map[string]EphemeralFieldMapping{
			"intelligence_config": {ApiName: "intelligenceConfig", AttrType: types.StringType},
			"effective_edition":   {ApiName: "effectiveEdition", AttrType: types.StringType},
		},
	}

	raw := map[string]interface{}{
		"intelligenceConfig": "projects/p/locations/global/intelligenceConfig",
		"effectiveEdition":   "STANDARD",
	}

	got, diags := FlattenEphemeralValue(raw, mapping)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	want, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"intelligence_config": types.StringValue("projects/p/locations/global/intelligenceConfig"),
		"effective_edition":   types.StringValue("STANDARD"),
	})
	if diags.HasError() {
		t.Fatalf("unexpected error building want value: %v", diags)
	}

	if !got.Equal(want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestFlattenEphemeralValue_ListOfObjects(t *testing.T) {
	elemAttrTypes := map[string]attr.Type{"item_name": types.StringType}
	elemType := types.ObjectType{AttrTypes: elemAttrTypes}
	mapping := EphemeralFieldMapping{
		AttrType: types.ListType{ElemType: elemType},
		Fields: map[string]EphemeralFieldMapping{
			"item_name": {ApiName: "itemName", AttrType: types.StringType},
		},
	}

	raw := []interface{}{
		map[string]interface{}{"itemName": "first"},
		map[string]interface{}{"itemName": "second"},
	}

	got, diags := FlattenEphemeralValue(raw, mapping)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	item1, _ := types.ObjectValue(elemAttrTypes, map[string]attr.Value{"item_name": types.StringValue("first")})
	item2, _ := types.ObjectValue(elemAttrTypes, map[string]attr.Value{"item_name": types.StringValue("second")})
	want, diags := types.ListValue(elemType, []attr.Value{item1, item2})
	if diags.HasError() {
		t.Fatalf("unexpected error building want value: %v", diags)
	}

	if !got.Equal(want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestFlattenEphemeralValue_MapOfObjects(t *testing.T) {
	elemAttrTypes := map[string]attr.Type{"value": types.StringType}
	elemType := types.ObjectType{AttrTypes: elemAttrTypes}
	mapping := EphemeralFieldMapping{
		AttrType: types.MapType{ElemType: elemType},
		Fields: map[string]EphemeralFieldMapping{
			"value": {ApiName: "value", AttrType: types.StringType},
		},
	}

	raw := map[string]interface{}{
		"key1": map[string]interface{}{"value": "map-item-1"},
		"key2": map[string]interface{}{"value": "map-item-2"},
	}

	got, diags := FlattenEphemeralValue(raw, mapping)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	v1, _ := types.ObjectValue(elemAttrTypes, map[string]attr.Value{"value": types.StringValue("map-item-1")})
	v2, _ := types.ObjectValue(elemAttrTypes, map[string]attr.Value{"value": types.StringValue("map-item-2")})
	want, diags := types.MapValue(elemType, map[string]attr.Value{"key1": v1, "key2": v2})
	if diags.HasError() {
		t.Fatalf("unexpected error building want value: %v", diags)
	}

	if !got.Equal(want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestFlattenEphemeralValue_ArbitraryDepthNesting(t *testing.T) {
	// Mirrors the Phase 0 2-level nesting spike: object-in-object, and
	// list-of-object-containing-object.
	innerAttrTypes := map[string]attr.Type{"inner_value": types.StringType}
	innerType := types.ObjectType{AttrTypes: innerAttrTypes}
	outerAttrTypes := map[string]attr.Type{
		"value": types.StringType,
		"inner": innerType,
	}
	outerType := types.ObjectType{AttrTypes: outerAttrTypes}

	innerFields := map[string]EphemeralFieldMapping{
		"inner_value": {ApiName: "innerValue", AttrType: types.StringType},
	}
	outerMapping := EphemeralFieldMapping{
		AttrType: outerType,
		Fields: map[string]EphemeralFieldMapping{
			"value": {ApiName: "value", AttrType: types.StringType},
			"inner": {ApiName: "inner", AttrType: innerType, Fields: innerFields},
		},
	}

	t.Run("object-in-object", func(t *testing.T) {
		raw := map[string]interface{}{
			"value": "outer-val",
			"inner": map[string]interface{}{"innerValue": "inner-val"},
		}
		got, diags := FlattenEphemeralValue(raw, outerMapping)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		innerVal, _ := types.ObjectValue(innerAttrTypes, map[string]attr.Value{"inner_value": types.StringValue("inner-val")})
		want, diags := types.ObjectValue(outerAttrTypes, map[string]attr.Value{
			"value": types.StringValue("outer-val"),
			"inner": innerVal,
		})
		if diags.HasError() {
			t.Fatalf("unexpected error building want value: %v", diags)
		}
		if !got.Equal(want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("list-of-object-containing-object", func(t *testing.T) {
		listMapping := EphemeralFieldMapping{
			AttrType: types.ListType{ElemType: outerType},
			Fields:   outerMapping.Fields,
		}

		raw := []interface{}{
			map[string]interface{}{
				"value": "item1",
				"inner": map[string]interface{}{"innerValue": "item1-inner"},
			},
		}

		got, diags := FlattenEphemeralValue(raw, listMapping)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags)
		}

		item1Inner, _ := types.ObjectValue(innerAttrTypes, map[string]attr.Value{"inner_value": types.StringValue("item1-inner")})
		item1, _ := types.ObjectValue(outerAttrTypes, map[string]attr.Value{
			"value": types.StringValue("item1"),
			"inner": item1Inner,
		})
		want, diags := types.ListValue(outerType, []attr.Value{item1})
		if diags.HasError() {
			t.Fatalf("unexpected error building want value: %v", diags)
		}
		if !got.Equal(want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestFlattenEphemeralValue_NullWhenApiFieldMissing(t *testing.T) {
	// Simulates res["effectiveIntelligenceConfig"] being absent from the API
	// response entirely (raw == nil, since res[missingKey] on a
	// map[string]interface{} returns the zero value).
	attrTypes := map[string]attr.Type{"intelligence_config": types.StringType}
	mapping := EphemeralFieldMapping{
		AttrType: types.ObjectType{AttrTypes: attrTypes},
		Fields: map[string]EphemeralFieldMapping{
			"intelligence_config": {ApiName: "intelligenceConfig", AttrType: types.StringType},
		},
	}

	got, diags := FlattenEphemeralValue(nil, mapping)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	want := types.ObjectNull(attrTypes)
	if !got.Equal(want) {
		t.Errorf("got %#v, want %#v (null object)", got, want)
	}
}
