package tpgresource

import (
	"reflect"
	"testing"
)

func TestSetOmittedFields(t *testing.T) {
	type Inner struct {
		InnerNotOmitted string   `json:"notOmitted"`
		InnerOmitted    []string `json:"-"`
	}
	type InputOuter struct {
		NotOmitted      string   `json:"notOmitted"`
		Omitted         []string `json:"-"`
		Struct          Inner
		Pointer         *Inner
		StructSlice     []Inner
		PointerSlice    []*Inner
		Unset           *Inner
		OnlyInInputType *Inner
	}
	type OutputOuter struct {
		NotOmitted       string   `json:"notOmitted"`
		Omitted          []string `json:"-"`
		Struct           Inner
		Pointer          *Inner
		StructSlice      []Inner
		PointerSlice     []*Inner
		Unset            *Inner
		OnlyInOutputType *Inner
	}

	input := &InputOuter{
		NotOmitted: "foo",
		Omitted:    []string{"foo"},
		Struct: Inner{
			InnerNotOmitted: "foo",
			InnerOmitted:    []string{"foo"},
		},
		Pointer: &Inner{
			InnerNotOmitted: "foo",
			InnerOmitted:    []string{"foo"},
		},
		StructSlice: []Inner{
			{
				InnerNotOmitted: "foo",
				InnerOmitted:    []string{"foo"},
			}, {
				InnerNotOmitted: "bar",
				InnerOmitted:    []string{"bar"},
			},
		},
		PointerSlice: []*Inner{
			{
				InnerNotOmitted: "foo",
				InnerOmitted:    []string{"foo"},
			}, {
				InnerNotOmitted: "bar",
				InnerOmitted:    []string{"bar"},
			},
		},
		OnlyInInputType: &Inner{
			InnerNotOmitted: "foo",
			InnerOmitted:    []string{"foo"},
		},
	}
	output := &OutputOuter{}
	if err := Convert(input, output); err != nil {
		t.Errorf("Error converting: %v", err)
	}
	if input.NotOmitted != output.NotOmitted ||
		!reflect.DeepEqual(input.Omitted, output.Omitted) ||
		!reflect.DeepEqual(input.Struct, output.Struct) ||
		!reflect.DeepEqual(input.Pointer, output.Pointer) ||
		!reflect.DeepEqual(input.StructSlice, output.StructSlice) ||
		!reflect.DeepEqual(input.PointerSlice, output.PointerSlice) ||
		!(input.Unset == nil && output.Unset == nil) {
		t.Errorf("Structs were not equivalent after conversion:\nInput:%#v\nOutput: %#v", input, output)
	}
}

func TestCanonicalizeIp(t *testing.T) {
	cases := map[string]struct {
		Input    interface{}
		Expected interface{}
	}{
		"nil": {
			Input:    nil,
			Expected: nil,
		},
		"empty string": {
			Input:    "",
			Expected: "",
		},
		"IPv4": {
			Input:    "10.0.0.1",
			Expected: "10.0.0.1",
		},
		"IPv6 full": {
			Input:    "2600:1900:4000:90ae:0:0:0:0",
			Expected: "2600:1900:4000:90ae::",
		},
		"IPv6 compressed": {
			Input:    "2600:1900:4000:90ae::",
			Expected: "2600:1900:4000:90ae::",
		},
		"non-ip string": {
			Input:    "not-an-ip",
			Expected: "not-an-ip",
		},
		"integer": {
			Input:    80,
			Expected: 80,
		},
	}

	for tn, tc := range cases {
		actual := CanonicalizeIp(tc.Input)
		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Errorf("%s failed: expected %#v, got %#v", tn, tc.Expected, actual)
		}
	}
}
