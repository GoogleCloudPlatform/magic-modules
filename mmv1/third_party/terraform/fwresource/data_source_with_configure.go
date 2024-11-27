package fwresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type dataSourceConfigurer interface {
	Configure(context.Context, datasource.ConfigureRequest, *datasource.ConfigureResponse)
}

type DataSourceWithConfigure struct {
	Config *transport_tpg.Config
}

func (ds *DataSourceWithConfigure) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	ds.Config = p
}
