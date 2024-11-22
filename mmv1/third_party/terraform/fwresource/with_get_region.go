// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwresource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type RegionGetter interface {
	GetRegion(types.String, *transport_tpg.Config, *diag.Diagnostics) types.String
}

// WithGetRegion can be embedded into _model) structs for resources and data sources
// that are expected to contain a Region field
type WithGetRegion struct {
	Config *transport_tpg.Config
}

// GetRegion determines whether to use a region value from the resource/data source config or the provider-default region
func (w *WithGetRegion) GetRegion(resourceRegion types.String, diags *diag.Diagnostics) types.String {
	// Prevent panic if nil pointer
	if w.Config == nil {
		diags.AddError("nil pointer encountered in fwresource.GetRegion", "Please report this issue to the provider developers")
		return types.String{}
	}

	if !(resourceRegion.IsNull() || resourceRegion.IsUnknown() || resourceRegion.ValueString() == "") {
		return resourceRegion
	}

	if w.Config.Region != "" {
		return types.StringValue(w.Config.Region)
	}

	diags.AddAttributeError(
		path.Root("region"),
		"missing required field region",
		"this resource/data source requires a region argument set either in the resource/data block or provided as a provider-level default",
	)
	return types.String{}
}
