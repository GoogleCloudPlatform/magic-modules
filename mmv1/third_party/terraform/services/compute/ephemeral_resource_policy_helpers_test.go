package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/fwutils"
)

func TestEphemeralDiskConsistencyGroupPolicy(t *testing.T) {
	attrs := map[string]attr.Type{"enabled": types.BoolType}
	mapping := fwutils.EphemeralFieldMapping{
		AttrType:  types.ObjectType{AttrTypes: attrs},
		Transform: flattenEphemeralDiskConsistencyGroupPolicy,
		Fields: map[string]fwutils.EphemeralFieldMapping{
			"enabled": {ApiName: "enabled", AttrType: types.BoolType},
		},
	}
	for _, raw := range []interface{}{nil, map[string]interface{}{}} {
		value, diags := fwutils.FlattenEphemeralValue(raw, mapping)
		if diags.HasError() {
			t.Fatal(diags)
		}
		if raw == nil {
			if !value.IsNull() {
				t.Fatal("absent policy must remain null")
			}
		} else if !value.(types.Object).Attributes()["enabled"].Equal(types.BoolValue(true)) {
			t.Fatal("present empty policy must be enabled")
		}
	}
}
