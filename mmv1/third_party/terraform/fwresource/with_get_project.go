package fwresource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type ProjectGetter interface {
	GetProject(types.String, *transport_tpg.Config, *diag.Diagnostics) types.String
}

type WithGetProject struct {
}

// GetProject determines whether to use a project value from the resource/data source config or the provider-default project
func (w *WithGetProject) GetProject(resourceProject types.String, config *transport_tpg.Config, diags *diag.Diagnostics) types.String {
	// Prevent panic if nil pointer
	if config == nil {
		diags.AddError("nil pointer encountered in fwresource.GetProject", "Please report this issue to the provider developers")
		return types.String{}
	}

	if !(resourceProject.IsNull() || resourceProject.IsUnknown() || resourceProject.ValueString() == "") {
		return resourceProject
	}

	if config.Project != "" {
		return types.StringValue(config.Project)
	}

	diags.AddAttributeError(
		path.Root("project"),
		"missing required field project",
		"this resource/data source requires a project argument set either in the resource/data block or provided as a provider-level default",
	)
	return types.String{}
}
