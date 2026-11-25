package biglakeiceberg

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBiglakeIcebergIcebergViewSchemaFieldsType(t *testing.T) {
	cases := map[string]struct {
		API    interface{}
		Config string
	}{
		"primitive type": {
			API:    "long",
			Config: "long",
		},
		"nested type": {
			API:    map[string]interface{}{"type": "list", "element-id": float64(3), "element": "string", "element-required": false},
			Config: `{"element":"string","element-id":3,"element-required":false,"type":"list"}`,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			if got := flattenBiglakeIcebergIcebergViewSchemaFieldsType(tc.API, nil, nil); got != tc.Config {
				t.Errorf("flatten: got %#v, want %q", got, tc.Config)
			}
			got, err := expandBiglakeIcebergIcebergViewSchemaFieldsType(tc.Config, nil, nil)
			if err != nil || !reflect.DeepEqual(got, tc.API) {
				t.Errorf("expand: got %#v (err: %v), want %#v", got, err, tc.API)
			}
		})
	}

	if _, err := expandBiglakeIcebergIcebergViewSchemaFieldsType(`{"type":"list"`, nil, nil); err == nil {
		t.Error("expand: expected an error for invalid JSON")
	}
}

func TestIcebergViewLocationDiffSuppress(t *testing.T) {
	cases := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{
			name: "unset in config",
			old:  "gs://bucket/ns/view/abc123def456ghij",
			new:  "",
			want: true,
		},
		{
			name: "exact match with trailing slash normalization",
			old:  "gs://bucket/ns/view",
			new:  "gs://bucket/ns/view/",
			want: true,
		},
		{
			name: "server random suffix under configured prefix",
			old:  "gs://bucket/ns/view/abc123def456ghij",
			new:  "gs://bucket/ns/view",
			want: true,
		},
		{
			name: "sibling prefix does not match",
			old:  "gs://bucket/ns/view_other/abc123def456ghij",
			new:  "gs://bucket/ns/view",
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := icebergViewLocationDiffSuppress("location", tc.old, tc.new, nil); got != tc.want {
				t.Errorf("icebergViewLocationDiffSuppress(%q, %q) = %v, want %v", tc.old, tc.new, got, tc.want)
			}
		})
	}
}

func TestIcebergViewPropertiesDiffSuppress(t *testing.T) {
	if !icebergViewPropertiesDiffSuppress("properties.write.metadata.compression-codec", "gzip", "", nil) {
		t.Error("expected server default compression-codec to be suppressed when unset")
	}
	if !icebergViewPropertiesDiffSuppress("properties.location", "gs://bucket/ns/view", "", nil) {
		t.Error("expected server default location property to be suppressed when unset")
	}
	if icebergViewPropertiesDiffSuppress("properties.write.metadata.compression-codec", "gzip", "none", nil) {
		t.Error("expected explicit compression-codec change not to be suppressed")
	}
	if icebergViewPropertiesDiffSuppress("properties.comment", "old", "", nil) {
		t.Error("expected user-defined property removal not to be suppressed")
	}
}

func TestIcebergViewJsonInt(t *testing.T) {
	if got, ok := icebergViewJsonInt(json.Number("42")); !ok || got != 42 {
		t.Errorf("json.Number: got (%d, %v), want (42, true)", got, ok)
	}
	if got, ok := icebergViewJsonInt(float64(7)); !ok || got != 7 {
		t.Errorf("float64: got (%d, %v), want (7, true)", got, ok)
	}
	if got, ok := icebergViewJsonInt(3); !ok || got != 3 {
		t.Errorf("int: got (%d, %v), want (3, true)", got, ok)
	}
	if _, ok := icebergViewJsonInt("not-an-int"); ok {
		t.Error("string: expected ok=false")
	}
}
