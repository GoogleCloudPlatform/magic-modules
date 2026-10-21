package fwutils

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandEphemeralValue(t *testing.T) {
	ctx := context.Background()
	list, _ := types.ListValue(types.StringType, []attr.Value{types.StringValue("first"), types.StringValue("second")})
	set, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("scope")})
	emptyList, _ := types.ListValue(types.StringType, []attr.Value{})
	emptyMap, _ := types.MapValue(types.StringType, map[string]attr.Value{})
	for _, tc := range []struct {
		name      string
		value     attr.Value
		field     EphemeralRequestField
		want      interface{}
		wantValue attr.Value
	}{
		{"false stays false", types.BoolValue(false), EphemeralRequestField{AttrType: types.BoolType, DefaultJSON: "true"}, false, types.BoolValue(false)},
		{"zero stays zero", types.Int64Value(0), EphemeralRequestField{AttrType: types.Int64Type, DefaultJSON: "42"}, int64(0), types.Int64Value(0)},
		{"empty string stays empty", types.StringValue(""), EphemeralRequestField{AttrType: types.StringType, DefaultJSON: `"default"`}, "", types.StringValue("")},
		{"double precision", types.Float64Value(1.23456789), EphemeralRequestField{AttrType: types.Float64Type}, 1.23456789, types.Float64Value(1.23456789)},
		{"list order", list, EphemeralRequestField{AttrType: list.Type(ctx)}, []interface{}{"first", "second"}, list},
		{"set", set, EphemeralRequestField{AttrType: set.Type(ctx)}, []interface{}{"scope"}, set},
		{"empty list", emptyList, EphemeralRequestField{AttrType: emptyList.Type(ctx)}, []interface{}{}, emptyList},
		{"empty map", emptyMap, EphemeralRequestField{AttrType: emptyMap.Type(ctx)}, map[string]interface{}{}, emptyMap},
		{"null stays null", types.StringNull(), EphemeralRequestField{AttrType: types.StringType}, nil, types.StringNull()},
		{"default on null", types.StringNull(), EphemeralRequestField{AttrType: types.StringType, DefaultJSON: `"3600s"`}, "3600s", types.StringValue("3600s")},
		{"false default", types.BoolNull(), EphemeralRequestField{AttrType: types.BoolType, DefaultJSON: "false"}, false, types.BoolValue(false)},
		{"large integer default", types.Int64Null(), EphemeralRequestField{AttrType: types.Int64Type, DefaultJSON: "9007199254740993"}, int64(9007199254740993), types.Int64Value(9007199254740993)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, value, diags := ExpandEphemeralValue(ctx, tc.value, tc.field)
			if diags.HasError() {
				t.Fatal(diags)
			}
			if !reflect.DeepEqual(body, tc.want) || !value.Equal(tc.wantValue) {
				t.Fatalf("got %#v / %v; want %#v / %v", body, value, tc.want, tc.wantValue)
			}
		})
	}
}

func TestExpandEphemeralNestedValues(t *testing.T) {
	ctx := context.Background()
	attrs := map[string]attr.Type{"display_name": types.StringType, "enabled": types.BoolType, "optional": types.StringType}
	object, _ := types.ObjectValue(attrs, map[string]attr.Value{"display_name": types.StringValue("example"), "enabled": types.BoolValue(false), "optional": types.StringNull()})
	child := EphemeralRequestField{AttrType: object.Type(ctx), Fields: map[string]EphemeralRequestField{
		"display_name": {ApiName: "displayName", AttrType: types.StringType, Required: true},
		"enabled":      {ApiName: "enabled", AttrType: types.BoolType},
		"optional":     {ApiName: "optional", AttrType: types.StringType},
	}}
	list, _ := types.ListValue(object.Type(ctx), []attr.Value{object})
	set, _ := types.SetValue(object.Type(ctx), []attr.Value{object})
	m, _ := types.MapValue(object.Type(ctx), map[string]attr.Value{"entry": object})
	for _, v := range []attr.Value{object, list, set, m} {
		field := child
		if !v.Type(ctx).Equal(object.Type(ctx)) {
			field = EphemeralRequestField{AttrType: v.Type(ctx), Element: &child}
		}
		body, result, diags := ExpandEphemeralValue(ctx, v, field)
		if diags.HasError() {
			t.Fatal(diags)
		}
		if !result.Equal(v) {
			t.Fatal("configured nested input changed")
		}
		wantObject := map[string]interface{}{"displayName": "example", "enabled": false}
		var want interface{} = wantObject
		switch v.(type) {
		case types.List, types.Set:
			want = []interface{}{wantObject}
		case types.Map:
			want = map[string]interface{}{"entry": wantObject}
		}
		if !reflect.DeepEqual(body, want) {
			t.Fatalf("got %#v, want %#v", body, want)
		}
	}
}

func TestExpandEphemeralInputValidation(t *testing.T) {
	empty, _ := types.SetValue(types.StringType, []attr.Value{})
	unknownChild, _ := types.ListValue(types.StringType, []attr.Value{types.StringUnknown()})
	for _, tc := range []struct {
		value attr.Value
		field EphemeralRequestField
	}{
		{types.StringUnknown(), EphemeralRequestField{AttrType: types.StringType}},
		{types.StringNull(), EphemeralRequestField{AttrType: types.StringType, Required: true}},
		{empty, EphemeralRequestField{AttrType: empty.Type(context.Background()), MinSize: 1}},
		{types.StringValue("invalid"), EphemeralRequestField{AttrType: types.StringType, Enum: []string{"valid"}}},
		{types.StringValue("-1"), EphemeralRequestField{AttrType: types.StringType, Pattern: `^[0-9]+s$`}},
		{unknownChild, EphemeralRequestField{AttrType: unknownChild.Type(context.Background())}},
	} {
		if _, _, diags := ExpandEphemeralValue(context.Background(), tc.value, tc.field); !diags.HasError() {
			t.Fatalf("expected input error for %v", tc.value)
		}
	}
}
