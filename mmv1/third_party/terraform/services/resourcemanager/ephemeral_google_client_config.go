package resourcemanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ ephemeral.EphemeralResource = &googleEphemeralClientConfig{}

func GoogleEphemeralClientConfig() ephemeral.EphemeralResource {
	return &googleEphemeralClientConfig{}
}

type googleEphemeralClientConfig struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralClientConfig) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_manager_secret_version"
}

type ephemeralClientConfigModel struct {
	Id            types.String `tfsdk:"id"`
	Project       types.String `tfsdk:"project"`
	Region        types.String `tfsdk:"region"`
	Zone          types.String `tfsdk:"zone"`
	AccessToken   types.String `tfsdk:"access_token"`
	DefaultLabels types.Map    `tfsdk:"default_labels"`
}

func (p *googleEphemeralClientConfig) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema.Description = "This ephemeral resource provides access to a Google Client config."
	resp.Schema.MarkdownDescription = "This ephemeral resource provides access to a Google Client config."

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of this data source in Terraform state. It is created in a projects/{{project}}/regions/{{region}}/zones/{{zone}} format and is NOT used by the data source in requests to Google APIs.",
				MarkdownDescription: "The ID of this data source in Terraform state. It is created in a projects/{{project}}/regions/{{region}}/zones/{{zone}} format and is NOT used by the data source in requests to Google APIs.",
			},
			"project": schema.StringAttribute{
				Description:         "The ID of the project to apply any resources to.",
				MarkdownDescription: "The ID of the project to apply any resources to.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				Description:         "The region to operate under.",
				MarkdownDescription: "The region to operate under.",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "The zone to operate under.",
				MarkdownDescription: "The zone to operate under.",
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				Description:         "The OAuth2 access token used by the client to authenticate against the Google Cloud API.",
				MarkdownDescription: "The OAuth2 access token used by the client to authenticate against the Google Cloud API.",
				Computed:            true,
				Sensitive:           true,
			},
			"default_labels": schema.MapAttribute{
				Description:         "The default labels configured on the provider.",
				MarkdownDescription: "The default labels configured on the provider.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (p *googleEphemeralClientConfig) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing userAgent and passing as an argument into a util function
	p.providerConfig = pd
}

func (p *googleEphemeralClientConfig) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data GoogleClientConfigModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("projects/%s/regions/%s/zones/%s", p.providerConfig.Project, p.providerConfig.Region, p.providerConfig.Zone))
	data.Project = types.StringValue(p.providerConfig.Project)
	data.Region = types.StringValue(p.providerConfig.Region)
	data.Zone = types.StringValue(p.providerConfig.Zone)

	// Convert default labels from SDK type system to plugin-framework data type
	m := map[string]*string{}
	for k, v := range p.providerConfig.DefaultLabels {
		// m[k] = types.StringValue(v)
		val := v
		m[k] = &val
	}
	dls, diags := types.MapValueFrom(ctx, types.StringType, m)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.DefaultLabels = dls

	token, err := p.providerConfig.TokenSource.Token()
	if err != nil {
		resp.Diagnostics.AddError("Error setting access_token", err.Error())
		return
	}
	data.AccessToken = types.StringValue(token.AccessToken)

	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}
