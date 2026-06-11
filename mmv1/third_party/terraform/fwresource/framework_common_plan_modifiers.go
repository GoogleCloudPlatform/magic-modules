package fwresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DefaultProjectModify(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, providerConfigProject string) {
	var old types.String
	diags := req.State.GetAttribute(ctx, path.Root("project"), &old)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var new types.String
	diags = req.Plan.GetAttribute(ctx, path.Root("project"), &new)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (old.IsUnknown() || old.IsNull()) && new.IsUnknown() {
		project := GetProjectFramework(new, types.StringValue(providerConfigProject), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Plan.SetAttribute(ctx, path.Root("project"), project)
	}
	return
}

func DefaultRegionModify(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, providerConfigRegion string) {
	var old types.String
	diags := req.State.GetAttribute(ctx, path.Root("region"), &old)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var new types.String
	diags = req.Plan.GetAttribute(ctx, path.Root("region"), &new)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (old.IsUnknown() || old.IsNull()) && new.IsUnknown() {
		region := GetRegionFramework(new, types.StringValue(providerConfigRegion), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Plan.SetAttribute(ctx, path.Root("region"), region)
	}
	return
}

func DefaultZoneModify(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, providerConfigZone string) {
	var old types.String
	diags := req.State.GetAttribute(ctx, path.Root("zone"), &old)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var new types.String
	diags = req.Plan.GetAttribute(ctx, path.Root("zone"), &new)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (old.IsUnknown() || old.IsNull()) && new.IsUnknown() {
		zone := GetZoneFramework(new, types.StringValue(providerConfigZone), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Plan.SetAttribute(ctx, path.Root("zone"), zone)
	}
	return
}
