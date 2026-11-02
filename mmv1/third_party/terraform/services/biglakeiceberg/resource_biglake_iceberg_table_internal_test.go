package biglakeiceberg

import (
	"reflect"
	"testing"
)

func TestBiglakeIcebergIcebergTableSchemaFieldsType(t *testing.T) {
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
			if got := flattenBiglakeIcebergIcebergTableSchemaFieldsType(tc.API, nil, nil); got != tc.Config {
				t.Errorf("flatten: got %#v, want %q", got, tc.Config)
			}
			got, err := expandBiglakeIcebergIcebergTableSchemaFieldsType(tc.Config, nil, nil)
			if err != nil || !reflect.DeepEqual(got, tc.API) {
				t.Errorf("expand: got %#v (err: %v), want %#v", got, err, tc.API)
			}
		})
	}

	if _, err := expandBiglakeIcebergIcebergTableSchemaFieldsType(`{"type":"list"`, nil, nil); err == nil {
		t.Error("expand: expected an error for invalid JSON")
	}
}
