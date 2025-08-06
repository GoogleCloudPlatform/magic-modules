package monitoring

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func FWMonitoringDashboardDiffSuppress() planmodifier.String {
	return &fwmonitoringDashboardDiffSuppress{}
}

type fwmonitoringDashboardDiffSuppress struct {
}

// Description returns a human-readable description of the plan modifier.
func (m fwmonitoringDashboardDiffSuppress) Description(_ context.Context) string {
	return "Verifies if computed attributes are the only difference in the dashboard_json field."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m fwmonitoringDashboardDiffSuppress) MarkdownDescription(_ context.Context) string {
	return "Verifies if computed attributes are the only difference in the dashboard_json field."
}

func (m *fwmonitoringDashboardDiffSuppress) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var old jsontypes.Normalized
	diags := req.State.GetAttribute(ctx, path.Root("dashboard_json"), &old)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var new jsontypes.Normalized
	diags = req.Plan.GetAttribute(ctx, path.Root("dashboard_json"), &new)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldMap, err := structure.ExpandJsonFromString(old.ValueString())
	if err != nil {
		return
	}
	newMap, err := structure.ExpandJsonFromString(new.ValueString())
	if err != nil {
		return
	}

	oldMap = recursiveRemoveComputedKeys(oldMap, newMap)

	if reflect.DeepEqual(oldMap, newMap) {
		resp.PlanValue = req.StateValue
	}

	return
}

// This recursive function takes an old map and a new map and is intended to remove the computed keys
// from the old json string (stored in state) so that it doesn't show a diff if it's not defined in the
// new map's json string (defined in config)
// this function is able to be modtly reused from the SDKv2 version of the resource
func recursiveRemoveComputedKeys(old map[string]interface{}, new map[string]interface{}) map[string]interface{} {
	for k, v := range old {
		if _, ok := old[k]; ok && new[k] == nil {
			delete(old, k)
			continue
		}

		if reflect.ValueOf(v).Kind() == reflect.Map {
			old[k] = removeComputedKeys(v.(map[string]interface{}), new[k].(map[string]interface{}))
			continue
		}

		if reflect.ValueOf(v).Kind() == reflect.Slice {
			for i, j := range v.([]interface{}) {
				if reflect.ValueOf(j).Kind() == reflect.Map && len(new[k].([]interface{})) > i {
					old[k].([]interface{})[i] = removeComputedKeys(j.(map[string]interface{}), new[k].([]interface{})[i].(map[string]interface{}))
				}
			}
			continue
		}
	}

	return old
}
