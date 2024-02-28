package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = RegionFromZoneFunction{}

func NewRegionFromZoneFunction() function.Function {
	return &RegionFromZoneFunction{}
}

type RegionFromZoneFunction struct{}

func (f RegionFromZoneFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "region_from_zone"
}

func (f RegionFromZoneFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the region within a provided resource's zone",
		Description: "Takes a single string argument, which should be a resource's zone.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "zone",
				Description: "A string of a resource's zone.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f RegionFromZoneFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var arg0 string
	resp.Diagnostics.Append(req.Arguments.GetArgument(ctx, 0, &arg0)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate input - can't be an empty string, every zone will have 2nd to last element in string be a `-`
	if len(arg0) != 0 && arg0[len(arg0)-2] == '-' {
		resp.Diagnostics.Append(resp.Result.Set(ctx, arg0[:len(arg0)-2])...)
	} else {
		return
	}

}
