package fwresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type resourceConfigurer interface {
	Configure(context.Context, resource.ConfigureRequest, *resource.ConfigureResponse)
}

type ResourceWithConfigure struct {
	Config *transport_tpg.Config
}

func (r *ResourceWithConfigure) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Config = p
}
