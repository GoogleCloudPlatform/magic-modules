package monitoring

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func MonitoringDashboardDiffSuppress() stringplanmodifier.String {
	return &monitoringDashboardDiffSuppress{}
}

type monitoringDashboardDiffSuppress struct {
}

func (d *instanceOptionDiffSuppress) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringRequest) {
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
		return false
	}
	newMap, err := structure.ExpandJsonFromString(new.ValueString())
	if err != nil {
		return false
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
